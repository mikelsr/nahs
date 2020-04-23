package events

import (
	"testing"

	"github.com/mikelsr/bspl"
)

func TestNewMessage(t *testing.T) {
	testNewMessageMarshal(t)
	testNewMessageUnmarshal(t)
}

func testNewMessageMarshal(t *testing.T) {
	i := testInstance()
	var m bspl.Message
	for _, v := range i.Messages() {
		if v.Action().Name == "Offer" {
			m = v
			break
		}
	}
	nm := MakeNewMessage(i.Key(), m)
	b, err := nm.Marshal()
	if err != nil {
		t.FailNow()
	}
	expectedLen := 335
	if len(b) != expectedLen {
		t.FailNow()
	}
}

func testNewMessageUnmarshal(t *testing.T) {
	i := testInstance()
	var m bspl.Message
	for _, v := range i.Messages() {
		if v.Action().Name == "Offer" {
			m = v
			break
		}
	}
	expected := MakeNewMessage(i.Key(), m)
	b, _ := expected.Marshal()
	nm, err := expected.Unmarshal(b)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	switch nm.Argument().(type) {
	case bspl.Message:
		break
	default:
		t.FailNow()
	}
	message := nm.Message()
	if !message.Parameters().Equals(m.Parameters()) ||
		message.Action().String() != m.Action().String() ||
		message.InstanceKey() != expected.InstanceKey() {
		t.FailNow()
	}
	if nm.ID() != expected.ID() || nm.Type() != expected.Type() {
		t.FailNow()
	}
}
