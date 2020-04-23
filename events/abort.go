package events

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
)

// Abort happens when a party cancels an Instance
type Abort struct {
	id          string
	instanceKey string
	motive      string
}

// MakeAbort is the default constructor for Abort
func MakeAbort(instanceKey string, motive string) Abort {
	return Abort{
		id:          uuid.New().String(),
		instanceKey: instanceKey,
		motive:      motive,
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

// InstanceKey returns the key of the instance of the Event
func (a Abort) InstanceKey() string {
	return a.instanceKey
}

// Marshal a Abort event to bytes
func (a Abort) Marshal() ([]byte, error) {
	motive := base64.StdEncoding.EncodeToString([]byte(a.motive))
	wrapper := EventWrapper{
		Argument:    string(motive),
		ID:          a.ID(),
		InstanceKey: a.instanceKey,
		Type:        TypeAbort,
	}
	return wrapper.Marshal()
}

// Motive for aborting the protocol
func (a Abort) Motive() string {
	return a.motive
}

// Unmarshal a Abort from bytes
func (a Abort) Unmarshal(data []byte) (Abort, error) {
	NIL := Abort{}
	wrapper := new(EventWrapper)
	if err := json.Unmarshal(data, wrapper); err != nil {
		return NIL, err
	}
	motive, err := base64.RawStdEncoding.DecodeString(wrapper.Argument)
	if err != nil {
		return NIL, err
	}
	n := Abort{
		id:          wrapper.ID,
		instanceKey: wrapper.InstanceKey,
		motive:      string(motive),
	}
	return n, nil
}
