package events

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mikelsr/bspl"
	imp "github.com/mikelsr/bspl/implementation"
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
	return ni.instance
}

// Type returns the event type
func (ni NewInstance) Type() EventType {
	return TypeNewInstance
}

// ID of the event
func (ni NewInstance) ID() string {
	return ni.id
}

// Instance returns the created instance
func (ni NewInstance) Instance() bspl.Instance {
	return ni.instance
}

// InstanceKey returns the key of the instance of the Event
func (ni NewInstance) InstanceKey() string {
	return ni.instance.Key()
}

// Marshal a NewInstance event to bytes
func (ni NewInstance) Marshal() ([]byte, error) {
	b, err := ni.instance.Marshal()
	if err != nil {
		return nil, err
	}
	instance := base64.StdEncoding.EncodeToString(b)
	wrapper := EventWrapper{
		Argument:    instance,
		ID:          ni.ID(),
		InstanceKey: ni.instance.Key(),
		Type:        TypeNewInstance,
	}
	return wrapper.Marshal()
}

// Unmarshal a NewInstance from bytes
func (ni NewInstance) Unmarshal(data []byte) (Event, error) {
	NIL := NewInstance{}
	wrapper := new(EventWrapper)
	if err := json.Unmarshal(data, wrapper); err != nil {
		return NIL, err
	}
	b, err := base64.StdEncoding.DecodeString(wrapper.Argument)
	if err != nil {
		return NIL, err
	}
	var instance imp.Instance
	dump, err := instance.Unmarshal(b)
	if err != nil {
		return NIL, err
	}
	instance = dump.(imp.Instance)
	n := NewInstance{
		id:       wrapper.ID,
		instance: instance,
	}
	return n, nil
}
