package net

import (
	"testing"

	log "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/crypto"
)

func TestMain(m *testing.M) {
	log.SetAllLoggers(log.LevelWarn)
	log.SetLogLevel(LogName, "info")

	testKeys = make([]*crypto.PrivKey, testNodeN)
	_loadTestKeys()
	_loadTestProtocols()
	m.Run()
}
