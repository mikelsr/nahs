package net

import (
	"context"
	"os"
	"path/filepath"

	"github.com/mikelsr/nahs/utils"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/pnet"
)

// Node represents a single NaHS peer.
type Node struct {
	// libp2p Host
	host host.Host
	// host cancelation function
	cancel context.CancelFunc
}

// NewNode is the default constructor for Node.
func NewNode(options ...libp2p.Option) *Node {
	n := new(Node)

	ctx, cancel := context.WithCancel(context.Background())
	n.cancel = cancel

	// Contatenate options parameter to default options
	opt := append(options, []libp2p.Option{
		// Attempt to open ports using uPNP for NATed hosts
		libp2p.NATPortMap(),
		// Let this host use relays and advertise itself on relays
		libp2p.EnableAutoRelay(),
	}...)

	// opt = append(opt, libp2p.PrivateNetwork(loadPrivNetPSK()))

	h, err := libp2p.New(ctx, opt...)
	if err != nil {
		panic(err)
	}
	n.host = h
	return n
}

// NodeFromPrivKey is a NewNode wrapper to create a new Node with the specified
// private key. Additional options may be provided.
func NodeFromPrivKey(sk crypto.PrivKey, options ...libp2p.Option) *Node {
	return NewNode(append(options, libp2p.Identity(sk))...)
}

// loadPrivNetPSK reads a private network PSK
func loadPrivNetPSK() pnet.PSK {
	dir, err := utils.GetProjectDir()
	if err != nil {
		panic(err)
	}
	file, err := os.Open(filepath.Join(dir, privNetPSKFile))
	if err != nil {
		panic(err)
	}
	psk, err := pnet.DecodeV1PSK(file)
	if err != nil {
		panic(err)
	}
	return psk
}
