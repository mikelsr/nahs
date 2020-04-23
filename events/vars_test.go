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

func (m mockReasoner) Abort(instanceKey string, motive string) error {
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

func (m mockReasoner) Instantiate(p bspl.Protocol, ins bspl.Values) (bspl.Instance, error) {
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
