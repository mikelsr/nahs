package events

import (
	"testing"

	"github.com/mikelsr/bspl"
)

func TestUpdateInstance(t *testing.T) {
	testUpdateEventMarshal(t)
	testUpdateEventUnmarshal(t)
}

func testUpdateEventMarshal(t *testing.T) {
	ue := MakeUpdateEvent(testInstance())
	b, err := ue.Marshal()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	expectedLen := 537
	if len(b) != expectedLen {
		t.FailNow()
	}
}

func testUpdateEventUnmarshal(t *testing.T) {
	expected := MakeUpdateEvent(testInstance())
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
	ue := event.(UpdateEvent)
	if ue.ID() != expected.ID() ||
		!ue.Instance().Equals(expected.instance) ||
		ue.Type() != expected.Type() ||
		ue.InstanceKey() != expected.InstanceKey() {
		t.FailNow()
	}
}
