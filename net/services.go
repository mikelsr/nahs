package net

import "github.com/mikelsr/bspl"

// AddProtocols adds a new BSPL Protocol to a Node
func (n *Node) AddProtocols(protocols ...bspl.Protocol) {
	for _, p := range protocols {
		n.protocols = append(n.protocols, p)
	}
}
