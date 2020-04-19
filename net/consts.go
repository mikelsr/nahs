package net

import (
	"path/filepath"

	log "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/protocol"
)

const (
	// LogName identifies the log of this module
	LogName = "nahs/net"

	listenAddrTCPIPv4 = "/ip4/0.0.0.0/tcp/0"
	listenAddrTCPIPv6 = "/ip6/::/tcp/0"

	// rendezvousString will identify the NaHS nodes at
	// the rendezvous points
	rendezvousString = "nahs-rendezvous"

	// ID of the BSPL discovery protocol
	protocolDiscoveryID = protocol.ID("/nahs/bspl/discovery/0.0.1")
)

var (
	logger = log.Logger(LogName)

	listenAddrs            = []string{listenAddrTCPIPv4, listenAddrTCPIPv6}
	privNetPSKFile         = filepath.Join("config", "private_network.psk")
	exchangeSeparator byte = '%'
	exchangeEnd       byte = '|'
)
