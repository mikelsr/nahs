package net

import (
	"context"

	"github.com/mikelsr/bspl"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

// Node represents a single NaHS peer.
type Node struct {
	// libp2p Host
	host host.Host
	// context of the node and the host
	context context.Context
	// host cancelation function
	cancel context.CancelFunc
	// dht table with information about network peers
	dht *dht.IpfsDHT
	// routing // TODO: rendezvous point?
	routing *discovery.RoutingDiscovery
	// services this node offers
	protocols []bspl.Protocol
}

// NewNode is the default constructor for Node.
func NewNode(options ...libp2p.Option) *Node {
	n := new(Node)

	n.context, n.cancel = context.WithCancel(context.Background())

	// Contatenate options parameter to default options
	opt := append(options, []libp2p.Option{
		libp2p.ListenAddrStrings(listenAddrs...),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let this host use relays and advertise itself on relays
		// libp2p.EnableAutoRelay(),
		// Attempt to open ports using uPNP for NATed hosts
		libp2p.NATPortMap(),
	}...)

	// opt = append(opt, libp2p.PrivateNetwork(loadPrivNetPSK()))

	h, err := libp2p.New(n.context, opt...)
	if err != nil {
		panic(err)
	}
	n.host = h

	// set stream handlers
	n.setStreamHandlers()

	logger.Infof("Created node with ID '%s'.", h.ID())
	return n
}

// NodeFromPrivKey is a NewNode wrapper to create a new Node with the specified
// private key. Additional options may be provided.
func NodeFromPrivKey(sk crypto.PrivKey, options ...libp2p.Option) *Node {
	return NewNode(append(options, libp2p.Identity(sk))...)
}

// ID of the libp2p host of the Node
func (n Node) ID() peer.ID {
	return n.host.ID()
}
