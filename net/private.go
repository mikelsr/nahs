package net

import (
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p-core/pnet"
	"github.com/mikelsr/nahs/utils"
)

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
