package net

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/mikelsr/bspl"
)

// Service contains a protocol describing the service and
// the roles the announcing node plays
type Service struct {
	Roles    []bspl.Role
	Protocol bspl.Protocol
}

// Services maps protocol keys to a Service with the protocol
// the key belongs to
type Services map[string]Service

// Contacts of peer: the nodes that have announced services
type Contacts map[peer.ID]Services

// AddContact adds a new contact to the Node
func (n *Node) AddContact(id peer.ID, services ...Service) {
	servs, found := n.Contacts[id]
	if !found {
		servs = make(Services)
		for _, s := range services {
			servs[s.Protocol.Key()] = s
		}
		n.Contacts[id] = servs
	}
	// it will react the same way and overwrite protocols with same id
	// but I'm leaving this branch here in case I change my mind
	if found {
		for _, s := range services {
			n.Contacts[id][s.Protocol.Key()] = s
		}
	}
}

// AddServices is functionally the same as AddContact for now.
func (n *Node) AddServices(id peer.ID, services ...Service) {
	n.AddContact(id, services...)
}

/*
func (n *Node) AddServices(id peer.ID, services Services) {
	servs := make([]Service, len(services))
	i := 0
	for _, v := range services {
		servs[i] = v
		i++
	}
	n.AddContact(id, servs...)
}
*/
