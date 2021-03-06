package net

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/mikelsr/nahs/events"
	"github.com/multiformats/go-multiaddr"
)

// setStreamHandler sets the stream handlers of the node peer
func (n *Node) setStreamHandlers() {
	n.host.SetStreamHandler(protocolDiscoveryID, n.discoveryHandler)
	n.host.SetStreamHandler(protocolEchoID, n.echoHandler)
	n.host.SetStreamHandler(protocolEventID, n.eventHandler)
}

func (n *Node) addRemotePeer(stream network.Stream) {
	// store new peer and multiaddr
	remotePeer := stream.Conn().RemotePeer()
	remoteAddrs := []multiaddr.Multiaddr{stream.Conn().RemoteMultiaddr()}
	logger.Debugf("Added address '%s' for peer '%s'", remoteAddrs[0], remotePeer.Pretty())
	n.host.Peerstore().AddAddrs(remotePeer, remoteAddrs, peerstore.PermanentAddrTTL)
}

// discoveryHandler exchanges the BSPL protocols of the
// services offered by each node
func (n *Node) discoveryHandler(stream network.Stream) {
	// defer recovery function in case the stream is closed
	// unexpectedly
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Recovered from error in protocol exchange: %s", r)
		}
		stream.Close()
	}()

	logger.Debug("Opened new BSPL protocol discovery stream")
	n.addRemotePeer(stream)

	var wg sync.WaitGroup
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	wg.Add(2)
	go n.discoveryReadData(rw, &wg, stream.Conn().RemotePeer())
	go n.discoveryWriteData(rw, &wg)
	wg.Wait()
}

// discoveryReadData parses the BSPL protocols transmitted by the other peer
func (n *Node) discoveryReadData(rw *bufio.ReadWriter, wg *sync.WaitGroup, sender peer.ID) {
	// defer recovery function in case the stream is closed
	// unexpectedly
	defer wg.Done()
	b, err := rw.ReadBytes(exchangeEnd)
	if err != nil {
		logger.Errorf("Error while reading protocol exchange: %s", err)
		panic(err)
	}
	bProtos := bytes.Split(b, []byte{exchangeSeparator})
	services := make([]Service, len(bProtos))

	// if  the protocol list was empty, return
	if len(bProtos) == 1 && len(bProtos[0]) == 1 && bytes.Equal(bProtos[0], []byte{exchangeEnd}) {
		logger.Debug("No new protocols discovered")
		return
	}
	// parse protocols
	for i, bp := range bProtos {
		protocol, roles, err := unwrapProtocol(bp[:len(bp)-1])
		if err != nil {
			panic(err)
		}
		services[i] = Service{
			Protocol: protocol,
			Roles:    roles,
		}
	}
	n.AddServices(sender, services...)
	var sb strings.Builder
	sb.WriteString("Discovered protocols: \n")
	for _, s := range services {
		sb.WriteString(s.Protocol.String())
	}
	logger.Debug(sb.String())
}

// discoveryWriteData transmits the BSPL protocols of this node to the other
func (n *Node) discoveryWriteData(rw *bufio.ReadWriter, wg *sync.WaitGroup) {
	defer wg.Done()
	k := len(n.protocols)
	for i, p := range n.protocols {
		roles := n.roles[p.Key()]
		if len(roles) == 0 {
			panic(fmt.Errorf("No defined roles for protocol '%s'", p.Key()))
		}
		payload := wrapProtocol(p, roles...)
		rw.Write(payload)
		if i != k-1 {
			rw.WriteByte(exchangeSeparator)
		}
	}
	rw.WriteByte(exchangeEnd)

	if err := rw.Flush(); err != nil {
		logger.Debugf("Error while writing protocol exchange: %s", err)
		panic(err)
	}
}

// echoHandler reads a message and writes the same as
// a response
func (n *Node) echoHandler(stream network.Stream) {
	// defer recovery function in case the stream is closed
	// unexpectedly
	defer func() {
		if r := recover(); r != nil {
			logger.Debugf("Recovered from error in protocol echo: %s", r)
		}
		stream.Close()

	}()
	logger.Debug("Opened new Echo stream")
	n.addRemotePeer(stream)

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	response := echoHandlerRead(rw)
	echoHandlerWrite(rw, response)

	stream.Close()
}

// echoHandlerRead and echoHandlerWrite are very short but useful for testing
func echoHandlerRead(rw *bufio.ReadWriter) []byte {
	b, err := rw.ReadBytes(exchangeEnd)
	if err != nil {
		logger.Debug("Error while reading echo message: %s", err)
		panic(err)
	}
	logger.Debugf("Received echo message: %s", string(b))
	return b
}

// echoHandlerWrite and echoHandlerRead are very short but useful for testing
func echoHandlerWrite(rw *bufio.ReadWriter, response []byte) {
	logger.Debugf("Send echo message: %s", string(response))
	rw.Write(response)
	if err := rw.Flush(); err != nil {
		logger.Debugf("Error while writing echo message: %s", err)
		panic(err)
	}
}

func (n *Node) eventHandler(stream network.Stream) {
	// defer recovery function in case the stream is closed
	// unexpectedly
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("Recovered from error in protocol event: %s", r)
			stream.Close()
		}
	}()
	logger.Debug("Opened new Echo stream")
	n.addRemotePeer(stream)
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	err := n.runEvent(rw, stream.Conn().RemotePeer())
	if err != nil {
		logger.Error(err)
		rw.Write(exchangeErr)
	} else {
		rw.Write(exchangeOk)
	}
	rw.WriteByte(exchangeEnd)
	if err := rw.Flush(); err != nil {
		logger.Debugf("Error while writing echo message: %s", err)
		panic(err)
	}
}

func (n *Node) runEvent(rw *bufio.ReadWriter, sender peer.ID) error {
	// read marshalled event
	b, err := rw.ReadBytes(exchangeEnd)
	b = b[:len(b)-1]
	if err != nil {
		logger.Error("Error while reading event message: %s", err)
		return err
	}
	// extract event ID
	id, err := events.ID(b)
	if err != nil {
		logger.Error(err)
		return ErrHandleEvent{ID: "-", Reason: "failed to extract ID"}
	}
	// err was already check with events.ID
	t, _ := events.Type(b)
	// extract instance key
	instanceKey, err := events.GetInstanceKey(b)
	if err != nil {
		logger.Error(err)
		return ErrHandleEvent{ID: id, Reason: "Could not extract instance key"}
	}
	// check if the instance has a peer assigned
	s, found := n.OpenInstances[instanceKey]
	logger.Debugf("Run event '%s' for node '%s'", t, sender)
	switch t {
	case events.TypeDropEvent, events.TypeUpdateEvent:
		// drop and update require an existing instance
		if !found {
			return ErrHandleEvent{ID: id, Reason: "Instance not found"}
		}
		// senders must coincide
		if s.String() != sender.String() {
			return ErrHandleEvent{ID: id, Reason: "Unauthorized"}
		}
		// remove event from OpenInstances
		if t == events.TypeDropEvent {
			delete(n.OpenInstances, instanceKey)
		}
	case events.TypeNewEvent:
		// new requires the instance to no exist
		if found {
			return ErrHandleEvent{ID: id, Reason: "Instance already existed"}
		}

		// asign sender to instance
		n.OpenInstances[instanceKey] = sender
	}
	// run event
	return events.RunEvent(n.reasoner, b)
}

func readEventResponse(rw *bufio.ReadWriter) (bool, error) {
	b, err := rw.ReadBytes(exchangeEnd)
	if err != nil {
		return false, err
	}
	if len(b) < len(exchangeOk) || len(b) < len(exchangeErr) {
		return false, errors.New("Response is too short")
	}
	// response is "ok"
	if bytes.Equal(b[:len(b)-1], exchangeOk) {
		return true, nil
	}
	return false, nil
}
