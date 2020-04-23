package events

import (
	"encoding/json"
	"errors"

	"github.com/mikelsr/bspl"
)

// EventType is used to differentiate events
type EventType string

const (
	// TypeAbort an instance was cancelled by one of the
	// parties involved
	TypeAbort EventType = "abort"
	// TypeNewInstance a new Instance was created
	TypeNewInstance EventType = "new_instance"
	// TypeNewMessage a new message was created
	// for an Instance
	TypeNewMessage EventType = "new_message"
)

// Event is used by one Node to notify another of a BSPL
// action
type Event interface {
	// Argument of an Event
	Argument() interface{}
	// ID of the Event
	ID() string
	// Instance Key
	InstanceKey() string
	// Type of Event
	Type() EventType
	// Marshal an Event to bytes
	Marshal() ([]byte, error)
	// Unmarshal an Event from bytes
	Unmarshal([]byte) (Event, error)
}

// EventWrapper is used by different event types
// to marshal themselves
type EventWrapper struct {
	Argument    string    `json:"arguments"`
	ID          string    `json:"id"`
	InstanceKey string    `json:"instance_key"`
	Type        EventType `json:"event_type"`
}

// Marshal an EventWrapper
func (e EventWrapper) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

type genericEvent struct {
	Type EventType `json:"event_type"`
	ID   string    `json:"id"`
}

// ID of the marshalled Event
func ID(marshalledEvent []byte) (string, error) {
	ge := new(genericEvent)
	if err := json.Unmarshal(marshalledEvent, ge); err != nil {
		return "", err
	}
	switch ge.Type {
	case TypeAbort, TypeNewInstance, TypeNewMessage:
		break
	default:
		return "", errors.New("Unable to identify event type")
	}
	return ge.ID, nil
}

// Type identities the type of a marshalled event
func Type(marshalledEvent []byte) (EventType, error) {
	ge := new(genericEvent)
	if err := json.Unmarshal(marshalledEvent, ge); err != nil {
		return "", err
	}
	switch ge.Type {
	case TypeAbort, TypeNewInstance, TypeNewMessage:
		break
	default:
		return "", errors.New("Unable to identify event type")
	}
	return ge.Type, nil
}

// GetInstanceKey extracts the instance key from a marshalled
// event
func GetInstanceKey(marshalledEvent []byte) (string, error) {
	t, err := Type(marshalledEvent)
	if err != nil {
		return "", err
	}
	var event Event
	switch t {
	case TypeAbort:
		var a Abort
		event, err = a.Unmarshal(marshalledEvent)
		if err != nil {
			return "", err
		}
	case TypeNewInstance:
		var ni NewInstance
		event, err = ni.Unmarshal(marshalledEvent)
		if err != nil {
			return "", err
		}
	case TypeNewMessage:
		var nm NewMessage
		event, err = nm.Unmarshal(marshalledEvent)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("Key not found")
	}
	return event.InstanceKey(), nil
}

// RunEvent identifies an event and calls the corresponding
// Reasoner method
func RunEvent(r bspl.Reasoner, marshalledEvent []byte) error {
	t, err := Type(marshalledEvent)
	if err != nil {
		return err
	}
	switch t {
	case TypeAbort:
		var a Abort
		event, err := a.Unmarshal(marshalledEvent)
		if err != nil {
			return err
		}
		a = event.(Abort)
		return r.Abort(a.InstanceKey(), a.Motive())
	case TypeNewInstance:
		var ni NewInstance
		event, err := ni.Unmarshal(marshalledEvent)
		if err != nil {
			return err
		}
		ni = event.(NewInstance)
		return r.RegisterInstance(ni.Instance())
	case TypeNewMessage:
		var nm NewMessage
		event, err := nm.Unmarshal(marshalledEvent)
		if err != nil {
			return err
		}
		nm = event.(NewMessage)
		return r.RegisterMessage(nm.InstanceKey(), nm.Message())
	}
	return nil
}
