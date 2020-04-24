package net

import (
	"bufio"
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/mikelsr/bspl"
	"github.com/mikelsr/nahs/events"
)

func TestDiscoveryHandler(t *testing.T) {
	n := testNodes(2)
	n1, n2 := n[0], n[1]

	// Add protocols to the nodes so they can advertise them
	n1.AddProtocol(tp1, tp1.Roles...)
	n2.AddProtocol(tp2, tp2.Roles...)

	// Create stream to exchange BSPL protocols
	stream, err := n1.host.NewStream(n1.context, n2.ID(), protocolDiscoveryID)
	if err != nil {
		n1.cancel()
		n2.cancel()
		t.Log(err)
		t.FailNow()
	}
	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	// Spawn and wait for RW routines
	var wg sync.WaitGroup
	wg.Add(2)
	go n1.discoveryReadData(rw, &wg, n2.ID())
	go n1.discoveryWriteData(rw, &wg)
	wg.Wait()
	if len(n1.Contacts) != 1 {
		t.FailNow()
	}
	for id, services := range n1.Contacts {
		if id != n2.ID() {
			t.FailNow()
		}
		if len(services) != 1 {
			t.FailNow()
		}
		for key, service := range services {
			if key != tp2.Key() || service.Protocol.String() != tp2.String() ||
				len(service.Roles) != 2 {
				t.FailNow()
			}
			break
		}
		break
	}
}

func TestEchoHandler(t *testing.T) {
	n := testNodes(2)
	n1, n2 := n[0], n[1]

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

func TestEventHandler(t *testing.T) {
	testEventHandlerAbort(t)
	testEventHandlerNewInstance(t)
	testEventHandlerNewMessage(t)
}

func testEventHandlerAbort(t *testing.T) {
	m := mockReasoner{}
	n := testNodes(3)
	for _, node := range n {
		node.reasoner = m
	}
	n1, n2, n3 := n[0], n[1], n[2]

	// create event
	instance := testInstance()

	n2.OpenInstances[instance.Key()] = n1.ID()

	a := events.MakeAbort(instance.Key(), "_")
	// send data from unauthorized node
	ok, err := n3.SendEvent(n2.ID(), a)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if ok {
		t.FailNow()
	}
	// send data from authorized node
	ok, err = n1.SendEvent(n2.ID(), a)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if !ok {
		t.FailNow()
	}
}

func testEventHandlerNewInstance(t *testing.T) {
	m := mockReasoner{}
	n := testNodes(2)
	for _, node := range n {
		node.reasoner = m
	}
	n1, n2 := n[0], n[1]

	// create event
	instance := testInstance()
	ni := events.MakeNewInstance(instance)

	// create new instance
	ok, err := n1.SendEvent(n2.ID(), ni)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if !ok {
		t.FailNow()
	}
	// create the same instance again
	_, err = n1.SendEvent(n2.ID(), ni)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
}

func testEventHandlerNewMessage(t *testing.T) {
	m := mockReasoner{}
	n := testNodes(2)
	for _, node := range n {
		node.reasoner = m
	}
	n1, n2 := n[0], n[1]

	// create event
	instance := testInstance()

	n2.OpenInstances[instance.Key()] = n1.ID()

	var message bspl.Message
	for _, v := range instance.Messages() {
		message = v
		break
	}

	nm := events.MakeNewMessage(instance.Key(), message)
	// send message to correct instance
	ok, err := n1.SendEvent(n2.ID(), nm)
	if err != nil {
		t.FailNow()
	}
	if !ok {
		t.FailNow()
	}
	// send message to incorrect instance
	nm = events.MakeNewMessage(instance.Key()+"_", message)
	ok, err = n1.SendEvent(n2.ID(), nm)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if ok {
		t.FailNow()
	}
}
