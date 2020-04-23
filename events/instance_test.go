package events

import (
	"testing"

	"github.com/mikelsr/bspl"
)

func TestNewIntance(t *testing.T) {
	testNewInstanceMarshal(t)
	testNewInstanceUnmarshal(t)
}

func testNewInstanceMarshal(t *testing.T) {
	ni := MakeNewInstance(testInstance())
	b, err := ni.Marshal()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	expectedLen := 1116
	if len(b) != expectedLen {
		t.FailNow()
	}
}

func testNewInstanceUnmarshal(t *testing.T) {
	expected := MakeNewInstance(testInstance())
	b, _ := expected.Marshal()
	ni, err := expected.Unmarshal(b)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	switch ni.Argument().(type) {
	case bspl.Instance:
		break
	default:
		t.FailNow()
	}
	if ni.ID() != expected.ID() ||
		!ni.Instance().Equals(expected.instance) ||
		ni.Type() != expected.Type() ||
		ni.InstanceKey() != expected.InstanceKey() {
		t.FailNow()
	}
}
