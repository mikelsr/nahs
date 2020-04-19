package net

import (
	"encoding/base64"
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/crypto"
)

// Storage types and functions for net components

// exportNode is a helper type to serialize Nodes
type exportNode struct {
	Key string `json:"prv"`
}

// Export a Node to bytes
func (n Node) Export() []byte {
	prv := n.host.Network().Peerstore().PrivKey(n.host.ID())
	b, _ := crypto.MarshalPrivateKey(prv)
	encoded := base64.StdEncoding.EncodeToString(b)
	e := exportNode{Key: encoded}
	result, _ := json.Marshal(e)
	return result
}
