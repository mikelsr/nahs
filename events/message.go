package events

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mikelsr/bspl"
)

// NewMessage happens when a party cancels an Instance
type NewMessage struct {
	id       string
	instance bspl.Instance
	message  bspl.Message
}

// MakeNewMessage is the default constructor for NewMessage
func MakeNewMessage(instance bspl.Instance, message bspl.Message) NewMessage {
	return NewMessage{
		id:       uuid.New().String(),
		instance: instance,
		message:  message,
	}
}

// Argument of NewMessage: nil.
func (nm NewMessage) Argument() interface{} {
	return nm.message
}

// Type returns the event type
func (nm NewMessage) Type() EventType {
	return TypeNewMessage
}

// ID of the event
func (nm NewMessage) ID() string {
	return nm.id
}

// Instance returns the instance of the Event
func (nm NewMessage) Instance() bspl.Instance {
	return nm.instance
}

// Marshal a NewMessage event to bytes
func (nm NewMessage) Marshal() ([]byte, error) {
	b, err := nm.instance.Marshal()
	if err != nil {
		return nil, err
	}
	instance := base64.StdEncoding.EncodeToString(b)
	b, err = nm.message.Marshal()
	if err != nil {
		return nil, err
	}
	message := base64.StdEncoding.EncodeToString(b)
	wrapper := EventWrapper{
		Argument: message,
		ID:       nm.ID(),
		Instance: instance,
		Type:     TypeNewMessage,
	}
	return wrapper.Marshal()
}

// Unmarshal a NewMessage from bytes
func (nm NewMessage) Unmarshal(data []byte) (NewMessage, error) {
	NIL := NewMessage{}
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
	b, err = base64.StdEncoding.DecodeString(wrapper.Argument)
	if err != nil {
		return NIL, err
	}
	var message bspl.Message
	message, err = message.Unmarshal(b)
	if err != nil {
		return NIL, err
	}
	n := NewMessage{
		id:       wrapper.ID,
		instance: instance,
		message:  message,
	}
	return n, nil
}
