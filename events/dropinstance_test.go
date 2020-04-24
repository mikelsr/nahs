package events

import (
	"testing"
)

func TestDropInstance(t *testing.T) {
	testDropInstanceMarshal(t)
	testDropInstanceUnmarshal(t)
}

func testDropInstanceMarshal(t *testing.T) {
	i := testInstance()
	motive := "the need to test this"
	a := MakeDropInstance(i.Key(), motive)
	b, err := a.Marshal()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	expectedLen := 148
	if len(b) != expectedLen {
		t.FailNow()
	}
}

func testDropInstanceUnmarshal(t *testing.T) {
	i := testInstance()
	motive := "the need to test this thing_"
	expected := MakeDropInstance(i.Key(), motive)
	b, _ := expected.Marshal()
	event, err := expected.Unmarshal(b)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	a := event.(DropInstance)
	switch a.Argument().(type) {
	case string:
		break
	default:
		t.FailNow()
	}
	if a.ID() != expected.ID() ||
		a.InstanceKey() != expected.InstanceKey() ||
		a.Motive() != expected.Motive() {
		t.FailNow()
	}
}
