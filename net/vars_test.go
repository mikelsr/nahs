package net

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p-core/crypto"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/mikelsr/bspl"
	imp "github.com/mikelsr/bspl/implementation"
	"github.com/mikelsr/bspl/proto"
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

	tp1, tp2 bspl.Protocol
	testKeys []*crypto.PrivKey
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

func testProtocol() proto.Protocol {
	buyer := proto.Role("Buyer")
	seller := proto.Role("Seller")
	p := proto.Protocol{
		Name:  "ProtoName",
		Roles: []proto.Role{buyer, seller},
		Params: []proto.Parameter{
			{Name: "ID", Key: true, Io: proto.Out},
			{Name: "item", Io: proto.Out},
			{Name: "price", Io: proto.Out},
		},
		Actions: []proto.Action{
			{Name: "Offer", From: buyer, To: seller, Params: []proto.Parameter{
				{Name: "ID", Key: true, Io: proto.In},
				{Name: "item", Io: proto.In},
				{Name: "price", Io: proto.Out},
			}},
			{Name: "Request", From: buyer, To: seller, Params: []proto.Parameter{
				{Name: "ID", Key: true, Io: proto.Out},
				{Name: "item", Io: proto.Out},
			}},
		},
	}
	return p
}

func testInstance() imp.Instance {
	p := testProtocol()
	roles := imp.Roles{
		proto.Role("Buyer"):  "B",
		proto.Role("Seller"): "S",
	}
	values := make(imp.Values)
	for _, param := range p.Parameters() {
		values[param.String()] = "X"
	}
	i := imp.NewInstance(p, roles, values)
	messages := make(imp.Messages)
	for _, a := range p.Actions {
		actionValues := make(imp.Values)
		for _, param := range a.Parameters() {
			actionValues[param.String()] = "X"
		}
		messages[a.String()] = imp.NewMessage(i.Key(), a, actionValues)
	}
	for _, m := range messages {
		i.AddMessage(m.(imp.Message))
	}
	return i
}

var (
	errMock error = errors.New("mock error")
)

type mockReasoner struct{}

func (m mockReasoner) DropInstance(instanceKey string, motive string) error {
	if instanceKey == testInstance().Key() {
		return nil
	}
	return errMock
}

func (m mockReasoner) GetInstance(instanceKey string) (bspl.Instance, bool) {
	if instanceKey == testInstance().Key() {
		return testInstance(), true
	}
	return imp.Instance{}, false
}
func (m mockReasoner) Instances(p bspl.Protocol) []bspl.Instance {
	return nil
}

func (m mockReasoner) Instantiate(p bspl.Protocol, roles bspl.Roles, ins bspl.Values) (bspl.Instance, error) {
	return nil, errMock
}

func (m mockReasoner) NewMessage(i bspl.Instance, a proto.Action) (bspl.Message, error) {
	return nil, errMock
}
func (m mockReasoner) RegisterInstance(i bspl.Instance) error {
	if i.Key() == testInstance().Key() {
		return nil
	}
	return errMock
}

func (m mockReasoner) RegisterMessage(instanceKey string, msg bspl.Message) error {
	if instanceKey == msg.InstanceKey() {
		return nil
	}
	return errMock
}

func testNodes(n int) []*Node {
	nodes := make([]*Node, n)
	for i := 0; i < n; i++ {
		nodes[i] = nodeFromPrivKey(*testKeys[i])
	}
	// Add addresses of each peer to the others
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			nodes[i].host.Peerstore().AddAddrs(nodes[j].ID(), nodes[j].Addrs(), peerstore.PermanentAddrTTL)
		}
	}
	return nodes
}
