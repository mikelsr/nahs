package events

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mikelsr/bspl"
	imp "github.com/mikelsr/bspl/implementation"
)

// NewEvent happens when an Instance is created
type NewEvent struct {
	id       string
	instance bspl.Instance
}

// MakeNewEvent is the default constructor for NewEvent
func MakeNewEvent(instance bspl.Instance) NewEvent {
	return NewEvent{
		id:       uuid.New().String(),
		instance: instance,
	}
}

// Argument of NewEvent: nil.
func (ne NewEvent) Argument() interface{} {
	return ne.instance
}

// Type returns the event type
func (ne NewEvent) Type() EventType {
	return TypeNewEvent
}

// ID of the event
func (ne NewEvent) ID() string {
	return ne.id
}

// Instance returns the created instance
func (ne NewEvent) Instance() bspl.Instance {
	return ne.instance
}

// InstanceKey returns the key of the instance of the Event
func (ne NewEvent) InstanceKey() string {
	return ne.instance.Key()
}

// Marshal a NewEvent event to bytes
func (ne NewEvent) Marshal() ([]byte, error) {
	b, err := ne.instance.Marshal()
	if err != nil {
		return nil, err
	}
	instance := base64.StdEncoding.EncodeToString(b)
	wrapper := EventWrapper{
		Argument:    instance,
		ID:          ne.ID(),
		InstanceKey: ne.instance.Key(),
		Type:        TypeNewEvent,
	}
	return wrapper.Marshal()
}

// Unmarshal a NewEvent from bytes
func (ne NewEvent) Unmarshal(data []byte) (Event, error) {
	NIL := NewEvent{}
	wrapper := new(EventWrapper)
	if err := json.Unmarshal(data, wrapper); err != nil {
		return NIL, err
	}
	b, err := base64.StdEncoding.DecodeString(wrapper.Argument)
	if err != nil {
		return NIL, err
	}
	instance := new(imp.Instance)

	if err = instance.Unmarshal(b); err != nil {
		return NIL, err
	}
	n := NewEvent{
		id:       wrapper.ID,
		instance: instance,
	}
	return n, nil
}
