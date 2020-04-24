package storage

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/mikelsr/bspl"
)

func loadTestHosts() ([]host.Host, []context.CancelFunc) {
	n := testNodeN
	hosts := make([]host.Host, n)
	cancels := make([]context.CancelFunc, n)
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
		ctx, cancel := context.WithCancel(context.Background())
		h, err := libp2p.New(ctx, libp2p.Identity(prv))
		if err != nil {
			panic(err)
		}
		hosts[i-1] = h
		cancels[i-1] = cancel
		if err != nil {
			panic(err)
		}
	}
	return hosts, cancels
}

func generateDB() {
	if _, err := os.Stat(testDBPath); !os.IsNotExist(err) {
		fmt.Println("Skipping test DB generation")
		return
	}
	fmt.Println("Generating test DB")
	hosts, cancels := loadTestHosts()
	r, err := os.Open(testBSPLFiles[0])
	if err != nil {
		panic(err)
	}
	p1, _ := bspl.Parse(r)
	r, err = os.Open(testBSPLFiles[1])
	if err != nil {
		panic(err)
	}
	p2, _ := bspl.Parse(r)

	k1 := hosts[0].ID()
	k2 := hosts[1].ID()

	v1 := []bspl.Protocol{p1, p2}
	v2 := []bspl.Protocol{p1}

	db, err := NewDB(testDBPath)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	db.Put(k1, v1)
	db.Put(k2, v2)
	for _, c := range cancels {
		c()
	}
}

func generatePeers() {
	n := 5
	var previousKey crypto.PrivKey
	for i := 1; i < n+1; i++ {
		prv, _, _ := crypto.GenerateRSAKeyPair(2048, rand.Reader)
		//id, _ := peer.IDFromPrivateKey(prv)
		//fmt.Printf("Generated ID for key %d: %s\n", i, id)
		if i == 1 {
			previousKey = prv
		} else {
			if previousKey.Equals(prv) {
				panic(errors.New("Equal keys"))
			}
		}
		//x, _ := prv.Raw()
		//fmt.Printf("Generated %d: %x\n", i, x)
		//b, err := crypto.MarshalPrivateKey(prv)
		b, err := crypto.MarshalPrivateKey(prv)
		if err != nil {
			panic(err)
		}
		encoded := base64.StdEncoding.EncodeToString(b)
		ioutil.WriteFile(filepath.Join(testKeysPath, fmt.Sprintf("peer_%d.key", i)), []byte(encoded), 0750)
	}
}
