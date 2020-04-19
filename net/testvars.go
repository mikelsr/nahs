package net

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/mikelsr/bspl"
	"github.com/mikelsr/nahs/utils"
)

var (
	testPath      = _genTestPath()
	testDBPath    = filepath.Join(testPath, "db", "test.db")
	testBSPLPath  = filepath.Join(testPath, "bspl")
	testBSPLFiles = []string{
		filepath.Join(testBSPLPath, "a.bspl"),
		filepath.Join(testBSPLPath, "x.bspl"),
	}
	testKeysPath = filepath.Join(testPath, "keys")
	testNodeN    = 5
)

func _genTestPath() string {
	dir, err := utils.GetProjectDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, "test")
}

func _loadTestProtocols() {
	// Load test BSPL protocols
	r, err := os.Open(testBSPLFiles[0])
	if err != nil {
		panic(err)
	}
	tp1, _ = bspl.Parse(r)
	r, err = os.Open(testBSPLFiles[1])
	if err != nil {
		panic(err)
	}
	tp2, _ = bspl.Parse(r)
}

func _loadTestKeys() {
	n := testNodeN

	for i := 1; i < n+1; i++ {
		b, err := ioutil.ReadFile(filepath.Join(
			testKeysPath,
			fmt.Sprintf("peer_%d.key", i)))
		if err != nil {
			panic(err)
		}
		decoded, err := base64.StdEncoding.DecodeString(string(b))
		if err != nil {
			panic(err)
		}
		prv, err := crypto.UnmarshalPrivateKey(decoded)
		if err != nil {
			panic(err)
		}
		testKeys[i-1] = &prv
	}
}
