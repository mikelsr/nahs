package events

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mikelsr/bspl"
	imp "github.com/mikelsr/bspl/implementation"
)

// UpdateEvent happens when an Instance is created
type UpdateEvent struct {
	id       string
	instance bspl.Instance
}

// MakeUpdateEvent is the default constructor for UpdateEvent
func MakeUpdateEvent(instance bspl.Instance) UpdateEvent {
	return UpdateEvent{
		id:       uuid.New().String(),
		instance: instance,
	}
}

// Argument of UpdateEvent: nil.
func (ue UpdateEvent) Argument() interface{} {
	return ue.instance
}

// Type returns the event type
func (ue UpdateEvent) Type() EventType {
	return TypeUpdateEvent
}

// ID of the event
func (ue UpdateEvent) ID() string {
	return ue.id
}

// Instance returns the created instance
func (ue UpdateEvent) Instance() bspl.Instance {
	return ue.instance
}

// InstanceKey returns the key of the instance of the Event
func (ue UpdateEvent) InstanceKey() string {
	return ue.instance.Key()
}

// Marshal a UpdateEvent event to bytes
func (ue UpdateEvent) Marshal() ([]byte, error) {
	b, err := ue.instance.Marshal()
	if err != nil {
		return nil, err
	}
	instance := base64.StdEncoding.EncodeToString(b)
	wrapper := EventWrapper{
		Argument:    instance,
		ID:          ue.ID(),
		InstanceKey: ue.instance.Key(),
		Type:        TypeUpdateEvent,
	}
	return wrapper.Marshal()
}

// Unmarshal a UpdateEvent from bytes
func (ue UpdateEvent) Unmarshal(data []byte) (Event, error) {
	NIL := UpdateEvent{}
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
	n := UpdateEvent{
		id:       wrapper.ID,
		instance: instance,
	}
	return n, nil
}
