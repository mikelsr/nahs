package events

import (
	"testing"
)

func TestDropEvent(t *testing.T) {
	testDropEventMarshal(t)
	testDropEventUnmarshal(t)
}

func testDropEventMarshal(t *testing.T) {
	i := testInstance()
	motive := "the need to test this"
	de := MakeDropEvent(i.Key(), motive)
	b, err := de.Marshal()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	expectedLen := 139
	if len(b) != expectedLen {
		t.FailNow()
	}
}

func testDropEventUnmarshal(t *testing.T) {
	i := testInstance()
	motive := "the need to test this thing_"
	expected := MakeDropEvent(i.Key(), motive)
	b, _ := expected.Marshal()
	event, err := expected.Unmarshal(b)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	de := event.(DropEvent)
	switch de.Argument().(type) {
	case string:
		break
	default:
		t.FailNow()
	}
	if de.ID() != expected.ID() ||
		de.InstanceKey() != expected.InstanceKey() ||
		de.Motive() != expected.Motive() {
		t.FailNow()
	}
}
