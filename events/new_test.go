package events

import (
	"testing"

	"github.com/mikelsr/bspl"
)

func TestNewInstance(t *testing.T) {
	testNewEventMarshal(t)
	testNewEventUnmarshal(t)
}

func testNewEventMarshal(t *testing.T) {
	ne := MakeNewEvent(testInstance())
	b, err := ne.Marshal()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	expectedLen := 534
	if len(b) != expectedLen {
		t.FailNow()
	}
}

func testNewEventUnmarshal(t *testing.T) {
	expected := MakeNewEvent(testInstance())
	b, _ := expected.Marshal()
	event, err := expected.Unmarshal(b)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	switch event.Argument().(type) {
	case bspl.Instance:
		break
	default:
		t.FailNow()
	}
	ne := event.(NewEvent)
	if ne.ID() != expected.ID() ||
		!ne.Instance().Equals(expected.instance) ||
		ne.Type() != expected.Type() ||
		ne.InstanceKey() != expected.InstanceKey() {
		t.FailNow()
	}
}
