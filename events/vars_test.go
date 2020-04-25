package events

import (
	"errors"

	"github.com/mikelsr/bspl"
	imp "github.com/mikelsr/bspl/implementation"
	"github.com/mikelsr/bspl/proto"
)

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

func testInstance() *imp.Instance {
	p := testProtocol()
	roles := imp.Roles{
		proto.Role("Buyer"):  "B",
		proto.Role("Seller"): "S",
	}
	i := imp.NewInstance(p, roles)
	i.SetValue("ID", "X")
	i.SetValue("item", "X")
	i.SetValue("price", "X")
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
	return nil, false
}

func (m mockReasoner) Instances(p bspl.Protocol) []bspl.Instance {
	return nil
}

func (m mockReasoner) Instantiate(p bspl.Protocol, roles bspl.Roles, ins bspl.Values) (bspl.Instance, error) {
	return nil, errMock
}

func (m mockReasoner) RegisterInstance(i bspl.Instance) error {
	if i.Key() == testInstance().Key() {
		return nil
	}
	return errMock
}

func (m mockReasoner) UpdateInstance(newVersion bspl.Instance) error {
	return nil
}
