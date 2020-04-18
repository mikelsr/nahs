package net

import (
	"path/filepath"

	"github.com/libp2p/go-libp2p-core/protocol"
)

const (
	listenAddrTCPIPv4 = "/ip4/0.0.0.0/tcp/0"
	listenAddrTCPIPv6 = "/ipv6/::/tcp/0"

	protocolID = protocol.ID("/nahs/bspl/0.0.1")
)

var (
	listenAddrs    = []string{listenAddrTCPIPv4, listenAddrTCPIPv6}
	privNetPSKFile = filepath.Join("config", "private_network.psk")
)
