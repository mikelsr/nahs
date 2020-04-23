package events

import (
	"testing"
)

func TestAbort(t *testing.T) {
	testAbortMarshal(t)
	testAbortUnmarshal(t)
}

func testAbortMarshal(t *testing.T) {
	i := testInstance()
	motive := "the need to test this"
	a := MakeAbort(i.Key(), motive)
	b, err := a.Marshal()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	expectedLen := 141
	if len(b) != expectedLen {
		t.FailNow()
	}
}

func testAbortUnmarshal(t *testing.T) {
	i := testInstance()
	motive := "the need to test this"
	expected := MakeAbort(i.Key(), motive)
	b, _ := expected.Marshal()
	a, err := expected.Unmarshal(b)
	if err != nil {
		t.FailNow()
	}
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
