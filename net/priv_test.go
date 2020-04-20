package net

import (
	"fmt"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
)

func TestPrivateNetwork(t *testing.T) {
	psk := loadPrivNetPSK()
	privNetOption := libp2p.PrivateNetwork(psk)
	// nodes 1 and 2 will belong to the private network
	// node 3 wont
	n1 := NodeFromPrivKey(*testKeys[0], privNetOption)
	n2 := NodeFromPrivKey(*testKeys[1], privNetOption)
	n3 := NodeFromPrivKey(*testKeys[2])
	// Add addresses of each peer to the others
	n1.host.Peerstore().AddAddrs(n2.ID(), n2.Addrs(), peerstore.PermanentAddrTTL)
	n1.host.Peerstore().AddAddrs(n3.ID(), n3.Addrs(), peerstore.PermanentAddrTTL)
	n2.host.Peerstore().AddAddrs(n3.ID(), n2.Addrs(), peerstore.PermanentAddrTTL)
	n2.host.Peerstore().AddAddrs(n3.ID(), n3.Addrs(), peerstore.PermanentAddrTTL)
	n3.host.Peerstore().AddAddrs(n1.ID(), n1.Addrs(), peerstore.PermanentAddrTTL)
	n3.host.Peerstore().AddAddrs(n2.ID(), n2.Addrs(), peerstore.PermanentAddrTTL)

	msg := []byte("Howdily neighbour")
	// Echo between two nodes in the same network
	if err := testEcho(n1, n2, msg); err != nil {
		n1.cancel()
		n2.cancel()
		t.FailNow()
	}

	// Test if n3 can stablish a connection with n1 without
	// timing out the test. If it can, the test will fail
	// because n3 is not in the private network
	timeout := time.After(1 * time.Second)
	done := make(chan bool)

	go func() {
		// Echo from a node from outside the network to a node inside the netwrok
		if err := testEcho(n3, n1, msg); err != nil {
			n1.cancel()
			n3.cancel()
			fmt.Println(err)
		}
		done <- true
	}()

	select {
	case <-timeout:
		n1.cancel()
		n2.cancel()
	case <-done:
		t.FailNow()
	}
}
