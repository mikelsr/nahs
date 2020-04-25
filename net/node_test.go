package net

import (
	"testing"

	log "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/mikelsr/bspl"
)

func TestMain(m *testing.M) {
	log.SetAllLoggers(log.LevelWarn)
	log.SetLogLevel(LogName, "debug")

	testKeys = make([]*crypto.PrivKey, testNodeN)
	_loadTestKeys()
	_loadTestProtocols()
	m.Run()
}

func TestNode_FindContact(t *testing.T) {
	m := mockReasoner{}
	p := testProtocol()
	n := testNodes(2)
	for _, node := range n {
		node.reasoner = m
	}
	n1, n2 := n[0], n[1]

	n1.AddContact(n2.ID(), Service{
		Roles:    []bspl.Role{"Seller"},
		Protocol: p,
	})

	contacts := n1.FindContact(p.Key(), "Seller")
	if len(contacts) != 1 || contacts[0] != n2.ID() {
		t.FailNow()
	}
}
