package net

import (
	"bufio"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
)

// configureDiscovery configures the node, connects to bootstrap nodes
// and announces self in the nodes
func (n *Node) configDiscovery() {
	// A local DHT will store network information in case bootstrap nodes
	// go down
	kademliaDHT, err := dht.New(n.context, n.host)
	if err != nil {
		panic(err)
	}

	// With the default configuration this will spawn a background thread
	// that will refresh the peer table ever 5 minutes
	logger.Debug("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(n.context); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	// Use default IPFS bootstrap peers
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := n.host.Connect(n.context, *peerinfo); err != nil {
				logger.Warning(err)
			} else {
				logger.Debug("Connection established with bootstrap node:", *peerinfo)
			}
		}()
	}
	wg.Wait()

	// Announce this node
	n.Announce()
}

// Announce self in network
func (n *Node) Announce() {
	logger.Debug("Announce self")
	routingDiscovery := discovery.NewRoutingDiscovery(n.dht)
	discovery.Advertise(n.context, routingDiscovery, rendezvousString)
}

// FindNodes searches for other NaHS nodes in the network
func (n *Node) FindNodes() {
	// Look for other NaHS nodes that have announced themselves
	logger.Debug("Search for other peers")
	peerChan, err := n.routing.FindPeers(n.context, rendezvousString)
	if err != nil {
		panic(err)
	}
	for peer := range peerChan {
		if peer.ID == n.ID() {
			continue
		}
		logger.Debugf("Found peer: %s", peer.ID)
		n.host.Peerstore().AddAddrs(peer.ID, peer.Addrs, peerstore.PermanentAddrTTL)

		// Exchange known services with the node
		stream, err := n.host.NewStream(n.context, peer.ID, protocolDiscoveryID)
		if err != nil {
			n.cancel()
			panic(err)
		}
		rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
		// Spawn routines to send the services offered by this node
		// to the other node
		// The wait group will be ignored
		var wg sync.WaitGroup
		wg.Add(2)
		go n.discoveryReadData(rw, &wg, peer.ID)
		go n.discoveryWriteData(rw, &wg)
		wg.Wait()
	}
	// block execution of this routine permantently
	// select {}
}
