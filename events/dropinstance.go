package events

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
)

// DropInstance happens when a party cancels an Instance
type DropInstance struct {
	id          string
	instanceKey string
	motive      string
}

// MakeDropInstance is the default constructor for DropInstance
func MakeDropInstance(instanceKey string, motive string) DropInstance {
	return DropInstance{
		id:          uuid.New().String(),
		instanceKey: instanceKey,
		motive:      motive,
	}
}

// Argument of DropInstance: nil.
func (a DropInstance) Argument() interface{} {
	return a.motive
}

// Type returns the event type
func (a DropInstance) Type() EventType {
	return TypeDropInstance
}

// ID of the event
func (a DropInstance) ID() string {
	return a.id
}

// InstanceKey returns the key of the instance of the Event
func (a DropInstance) InstanceKey() string {
	return a.instanceKey
}

// Marshal a DropInstance event to bytes
func (a DropInstance) Marshal() ([]byte, error) {
	motive := base64.StdEncoding.EncodeToString([]byte(a.motive))
	wrapper := EventWrapper{
		Argument:    string(motive),
		ID:          a.ID(),
		InstanceKey: a.instanceKey,
		Type:        TypeDropInstance,
	}
	return wrapper.Marshal()
}

// Motive for dropInstanceing the protocol
func (a DropInstance) Motive() string {
	return a.motive
}

// Unmarshal a DropInstance from bytes
func (a DropInstance) Unmarshal(data []byte) (Event, error) {
	NIL := DropInstance{}
	wrapper := new(EventWrapper)
	if err := json.Unmarshal(data, wrapper); err != nil {
		return NIL, err
	}
	motive, err := base64.StdEncoding.DecodeString(string(wrapper.Argument))
	if err != nil {
		return NIL, err
	}
	n := DropInstance{
		id:          wrapper.ID,
		instanceKey: wrapper.InstanceKey,
		motive:      string(motive),
	}
	return n, nil
}
