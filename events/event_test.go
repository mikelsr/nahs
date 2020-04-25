package events

import (
	"testing"
)

func TestType(t *testing.T) {
	i := testInstance()
	motive := "general test"
	aEvent, _ := MakeDropEvent(i.Key(), motive).Marshal()
	niEvent, _ := MakeNewEvent(i).Marshal()
	nmEvent, _ := MakeUpdateEvent(i).Marshal()

	dumps := [][]byte{aEvent, niEvent, nmEvent}
	types := make([]EventType, 3)
	for j := 0; j < 3; j++ {
		x, err := Type(dumps[j])
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		types[j] = x
	}
	if types[0] != TypeDropEvent || types[1] != TypeNewEvent || types[2] != TypeUpdateEvent {
		t.FailNow()
	}
}

func TestRunEvent(t *testing.T) {
	i := testInstance()
	motive := "general test"
	dropEvent, _ := MakeDropEvent(i.Key(), motive).Marshal()
	newEvent, _ := MakeNewEvent(i).Marshal()
	updateEvent, _ := MakeUpdateEvent(i).Marshal()

	r := mockReasoner{}

	for _, event := range [][]byte{dropEvent, newEvent, updateEvent} {
		if err := RunEvent(r, event); err != nil {
			t.Log(err)
			t.FailNow()
		}
	}
	if err := RunEvent(r, []byte{}); err == nil {
		t.FailNow()
	}
}

func TestGetInstanceKey(t *testing.T) {
	i := testInstance()
	motive := "general test"
	aEvent, _ := MakeDropEvent(i.Key(), motive).Marshal()
	niEvent, _ := MakeNewEvent(i).Marshal()
	nmEvent, _ := MakeUpdateEvent(i).Marshal()

	for _, event := range [][]byte{aEvent, niEvent, nmEvent} {
		key, err := GetInstanceKey(event)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
		if key != i.Key() {
			t.FailNow()
		}
	}
	if _, err := GetInstanceKey([]byte{}); err == nil {
		t.FailNow()
	}
}
