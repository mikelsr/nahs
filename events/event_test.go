package events

import (
	"testing"

	"github.com/mikelsr/bspl"
)

func TestType(t *testing.T) {
	i := testInstance()
	motive := "general test"
	var m bspl.Message
	for _, v := range i.Messages() {
		if v.Action().Name == "Offer" {
			m = v
			break
		}
	}
	aEvent, _ := MakeDropInstance(i.Key(), motive).Marshal()
	niEvent, _ := MakeNewInstance(i).Marshal()
	nmEvent, _ := MakeNewMessage(i.Key(), m).Marshal()

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
	if types[0] != TypeDropInstance || types[1] != TypeNewInstance || types[2] != TypeNewMessage {
		t.FailNow()
	}
}

func TestRunEvent(t *testing.T) {
	i := testInstance()
	motive := "general test"
	var m bspl.Message
	for _, v := range i.Messages() {
		if v.Action().Name == "Offer" {
			m = v
			break
		}
	}
	aEvent, _ := MakeDropInstance(i.Key(), motive).Marshal()
	niEvent, _ := MakeNewInstance(i).Marshal()
	nmEvent, _ := MakeNewMessage(i.Key(), m).Marshal()

	r := mockReasoner{}

	for _, event := range [][]byte{aEvent, niEvent, nmEvent} {
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
	var m bspl.Message
	for _, v := range i.Messages() {
		if v.Action().Name == "Offer" {
			m = v
			break
		}
	}
	aEvent, _ := MakeDropInstance(i.Key(), motive).Marshal()
	niEvent, _ := MakeNewInstance(i).Marshal()
	nmEvent, _ := MakeNewMessage(i.Key(), m).Marshal()

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
