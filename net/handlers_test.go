package net

import (
	"bufio"
	"bytes"
	"fmt"
	"sync"
	"testing"

	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/mikelsr/bspl"
)

func TestDiscoveryHandler(t *testing.T) {
	n1 := NodeFromPrivKey(*testKeys[0])
	n2 := NodeFromPrivKey(*testKeys[1])
	// Add addresses of each peer to the others
	n1.host.Peerstore().AddAddrs(n2.ID(), n2.Addrs(), peerstore.PermanentAddrTTL)
	n2.host.Peerstore().AddAddrs(n1.ID(), n1.Addrs(), peerstore.PermanentAddrTTL)

	// Add protocols to the nodes so they can advertise them
	n1.protocols = []bspl.Protocol{tp1}
	n2.protocols = []bspl.Protocol{tp2}

	// Create stream to exchange BSPL protocols
	stream, err := n1.host.NewStream(n1.context, n2.ID(), protocolDiscoveryID)
	if err != nil {
		n1.cancel()
		n2.cancel()
		fmt.Println(err)
		t.FailNow()
	}
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	// Spawn and wait for RW routines
	var wg sync.WaitGroup
	wg.Add(2)
	go n1.discoveryReadData(rw, &wg)
	go n1.discoveryWriteData(rw, &wg)
	wg.Wait()
}

func TestEchoHandler(t *testing.T) {
	n1 := NodeFromPrivKey(*testKeys[0])
	n2 := NodeFromPrivKey(*testKeys[1])
	// Add addresses of each peer to the others
	n1.host.Peerstore().AddAddrs(n2.ID(), n2.Addrs(), peerstore.PermanentAddrTTL)
	n2.host.Peerstore().AddAddrs(n1.ID(), n1.Addrs(), peerstore.PermanentAddrTTL)
	// call testEcho and fail tests if it errs
	msg := []byte("Howdily doodily")
	if err := testEcho(n1, n2, msg); err != nil {
		n1.cancel()
		n2.cancel()
		fmt.Println(err)
		t.FailNow()
	}

}

func testEcho(n1, n2 *Node, msg []byte) error {
	// Create echo stream
	stream, err := n1.host.NewStream(n1.context, n2.ID(), protocolEchoID)
	if err != nil {
		return nil
	}
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	message := append(msg, exchangeEnd)
	rw.Write(message)
	if err := rw.Flush(); err != nil {
		return err
	}
	// Launch RW functions on order
	// Test will fail if it times out
	response := echoHandlerRead(rw)
	if !bytes.Equal(response, message) {
		return fmt.Errorf("Echo expected '%s' but got '%s'", message, response)
	}
	echoHandlerWrite(rw, response)
	return nil
}
