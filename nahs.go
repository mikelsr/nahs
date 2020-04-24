package nahs

import (
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/mikelsr/bspl"
	"github.com/mikelsr/nahs/net"
)

type (
	// Node of the NaHS network.
	Node = net.Node
)

// NewNode creates a new NaHS node. LibP2P options can be passed
// to configure the node.
func NewNode(reasoner bspl.Reasoner, options ...libp2p.Option) *Node {
	return net.NewNode(reasoner, options...)
}

// MakeNode creates a node with the specified private key so the
// node maintains the ID it previously had.
func MakeNode(reasoner bspl.Reasoner, sk crypto.PrivKey, options ...libp2p.Option) *Node {
	return net.NodeFromPrivKey(reasoner, sk, options...)
}
