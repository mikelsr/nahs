package net

import "path/filepath"

const (
	listenAddrTCPIPv4 = "/ip4/0.0.0.0/tcp/0"
	listenAddrTCPIPv6 = "/ipv6/::/tcp/0"
)

var (
	listenAddrs    = []string{listenAddrTCPIPv4, listenAddrTCPIPv6}
	privNetPSKFile = filepath.Join("config", "private_network.psk")
)
