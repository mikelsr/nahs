package net

import (
	"bufio"
	"bytes"
	"strings"
	"sync"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/mikelsr/bspl"
)

// setStreamHandler sets the stream handlers of the node peer
func (n *Node) setStreamHandlers() {
	n.host.SetStreamHandler(protocolDiscoveryID, n.discoveryHandler)
}

// discoveryHandler exchanges the BSPL protocols of the
// services offered by each node
func (n *Node) discoveryHandler(stream network.Stream) {
	// defer recovery function in case the stream is closed
	// unexpectedly
	defer func() {
		if r := recover(); r != nil {
			logger.Infof("Recovered from error in protocol exchange %s", r)
		}
	}()

	logger.Info("Opened new BSPL protocol discovery stream")

	var wg sync.WaitGroup
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	wg.Add(2)
	go n.discoveryReadData(rw, &wg)
	go n.discoveryWriteData(rw, &wg)
	wg.Wait()
	stream.Close()
}

// discoveryReadData parses
func (n *Node) discoveryReadData(rw *bufio.ReadWriter, wg *sync.WaitGroup) {
	// defer recovery function in case the stream is closed
	// unexpectedly
	defer wg.Done()
	b, err := rw.ReadBytes(exchangeEnd)
	if err != nil {
		logger.Info("Error while reading protocol exchange: %s", err)
		panic(err)
	}
	bProtos := bytes.Split(b, []byte{exchangeSeparator})
	protocols := make([]bspl.Protocol, len(bProtos))

	// if  the protocol list was empty, return
	if len(bProtos) == 1 && len(bProtos[0]) == 1 && bytes.Equal(bProtos[0], []byte{exchangeEnd}) {
		logger.Info("No new protocols discovered")
		return
	}
	// parse protocols
	for i, bp := range bProtos {
		reader := bytes.NewReader(bp)
		protocol, err := bspl.Parse(reader)
		if err != nil {
			panic(err)
		}
		protocols[i] = protocol
	}
	var sb strings.Builder
	sb.WriteString("Discovered protocols: \n")
	for _, p := range protocols {
		sb.WriteString(p.String())
	}
	logger.Info(sb.String())
}

func (n *Node) discoveryWriteData(rw *bufio.ReadWriter, wg *sync.WaitGroup) {
	defer wg.Done()
	k := len(n.protocols)
	for i, p := range n.protocols {
		rw.WriteString(p.String())
		if i != k-1 {
			rw.WriteByte(exchangeSeparator)
		}
	}
	rw.WriteByte(exchangeEnd)

	if err := rw.Flush(); err != nil {
		logger.Infof("Error while writing protocol exchange: %s", err)
		panic(err)
	}
}
