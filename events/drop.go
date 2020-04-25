package events

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
)

// DropEvent happens when de party cancels an Instance
type DropEvent struct {
	id          string
	instanceKey string
	motive      string
}

// MakeDropEvent is the default constructor for DropEvent
func MakeDropEvent(instanceKey string, motive string) DropEvent {
	return DropEvent{
		id:          uuid.New().String(),
		instanceKey: instanceKey,
		motive:      motive,
	}
}

// Argument of DropEvent: nil.
func (de DropEvent) Argument() interface{} {
	return de.motive
}

// Type returns the event type
func (de DropEvent) Type() EventType {
	return TypeDropEvent
}

// ID of the event
func (de DropEvent) ID() string {
	return de.id
}

// InstanceKey returns the key of the instance of the Event
func (de DropEvent) InstanceKey() string {
	return de.instanceKey
}

// Marshal de DropEvent event to bytes
func (de DropEvent) Marshal() ([]byte, error) {
	motive := base64.StdEncoding.EncodeToString([]byte(de.motive))
	wrapper := EventWrapper{
		Argument:    string(motive),
		ID:          de.ID(),
		InstanceKey: de.instanceKey,
		Type:        TypeDropEvent,
	}
	return wrapper.Marshal()
}

// Motive for dropInstanceing the protocol
func (de DropEvent) Motive() string {
	return de.motive
}

// Unmarshal de DropEvent from bytes
func (de DropEvent) Unmarshal(data []byte) (Event, error) {
	NIL := DropEvent{}
	wrapper := new(EventWrapper)
	if err := json.Unmarshal(data, wrapper); err != nil {
		return NIL, err
	}
	motive, err := base64.StdEncoding.DecodeString(string(wrapper.Argument))
	if err != nil {
		return NIL, err
	}
	n := DropEvent{
		id:          wrapper.ID,
		instanceKey: wrapper.InstanceKey,
		motive:      string(motive),
	}
	return n, nil
}
