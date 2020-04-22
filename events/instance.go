package events

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mikelsr/bspl"
)

// NewInstance happens when an Instance is created
type NewInstance struct {
	id       string
	instance bspl.Instance
}

// MakeNewInstance is the default constructor for NewInstance
func MakeNewInstance(instance bspl.Instance) NewInstance {
	return NewInstance{
		id:       uuid.New().String(),
		instance: instance,
	}
}

// Argument of NewInstance: nil.
func (ni NewInstance) Argument() interface{} {
	return nil
}

// Type returns the event type
func (ni NewInstance) Type() EventType {
	return TypeNewInstance
}

// ID of the event
func (ni NewInstance) ID() string {
	return ni.id
}

// Instance returns the instance of the Event
func (ni NewInstance) Instance() bspl.Instance {
	return ni.instance
}

// Marshal a NewInstance event to bytes
func (ni NewInstance) Marshal() ([]byte, error) {
	b, err := ni.instance.Marshal()
	if err != nil {
		return nil, err
	}
	instance := base64.StdEncoding.EncodeToString(b)
	wrapper := EventWrapper{
		Argument: "",
		ID:       ni.ID(),
		Instance: instance,
		Type:     TypeNewInstance,
	}
	return wrapper.Marshal()
}

// Unmarshal a NewInstance from bytes
func (ni NewInstance) Unmarshal(data []byte) (NewInstance, error) {
	NIL := NewInstance{}
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
	n := NewInstance{
		id:       wrapper.ID,
		instance: instance,
	}
	return n, nil
}
