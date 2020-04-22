package events

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mikelsr/bspl"
)

// Abort happens when a party cancels an Instance
type Abort struct {
	id       string
	instance bspl.Instance
	motive   string
}

// MakeAbort is the default constructor for Abort
func MakeAbort(instance bspl.Instance, motive string) Abort {
	return Abort{
		id:       uuid.New().String(),
		instance: instance,
		motive:   motive,
	}
}

// Argument of Abort: nil.
func (a Abort) Argument() interface{} {
	return a.motive
}

// Type returns the event type
func (a Abort) Type() EventType {
	return TypeAbort
}

// ID of the event
func (a Abort) ID() string {
	return a.id
}

// Instance returns the instance of the Event
func (a Abort) Instance() bspl.Instance {
	return a.instance
}

// Marshal a Abort event to bytes
func (a Abort) Marshal() ([]byte, error) {
	b, err := a.instance.Marshal()
	if err != nil {
		return nil, err
	}
	instance := base64.StdEncoding.EncodeToString(b)
	wrapper := EventWrapper{
		Argument: a.motive,
		ID:       a.ID(),
		Instance: instance,
		Type:     TypeAbort,
	}
	return wrapper.Marshal()
}

// Unmarshal a Abort from bytes
func (a Abort) Unmarshal(data []byte) (Abort, error) {
	NIL := Abort{}
	wrapper := new(EventWrapper)
	if err := json.Unmarshal(data, wrapper); err != nil {
		return NIL, err
	}
	b, err := base64.StdEncoding.DecodeString(wrapper.Instance)
	if err != nil {
		return NIL, err
	}
	var instance bspl.Instance
	instance, err = instance.Unmarshal(b)
	if err != nil {
		return NIL, err
	}
	n := Abort{
		id:       wrapper.ID,
		instance: instance,
		motive:   wrapper.Argument,
	}
	return n, nil
}
