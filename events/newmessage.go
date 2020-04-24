package events

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/mikelsr/bspl"
	imp "github.com/mikelsr/bspl/implementation"
)

// NewMessage happens when a party cancels an Instance
type NewMessage struct {
	id          string
	instanceKey string
	message     bspl.Message
}

// MakeNewMessage is the default constructor for NewMessage
func MakeNewMessage(instanceKey string, message bspl.Message) NewMessage {
	return NewMessage{
		id:          uuid.New().String(),
		instanceKey: instanceKey,
		message:     message,
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

// InstanceKey returns the key of the instance of the Event
func (nm NewMessage) InstanceKey() string {
	return nm.instanceKey
}

// Marshal a NewMessage event to bytes
func (nm NewMessage) Marshal() ([]byte, error) {
	b, err := nm.message.Marshal()
	if err != nil {
		return nil, err
	}
	message := base64.StdEncoding.EncodeToString(b)
	wrapper := EventWrapper{
		Argument:    message,
		ID:          nm.ID(),
		InstanceKey: nm.instanceKey,
		Type:        TypeNewMessage,
	}
	return wrapper.Marshal()
}

// Message returns the created Message
func (nm NewMessage) Message() bspl.Message {
	return nm.message
}

// Unmarshal a NewMessage from bytes
func (nm NewMessage) Unmarshal(data []byte) (Event, error) {
	NIL := NewMessage{}
	wrapper := new(EventWrapper)
	if err := json.Unmarshal(data, wrapper); err != nil {
		return NIL, err
	}
	b, err := base64.StdEncoding.DecodeString(wrapper.Argument)
	if err != nil {
		return NIL, err
	}
	var message imp.Message
	dump, err := message.Unmarshal(b)
	if err != nil {
		return NIL, err
	}
	n := NewMessage{
		id:          wrapper.ID,
		instanceKey: wrapper.InstanceKey,
		message:     dump.(imp.Message),
	}
	return n, nil
}
