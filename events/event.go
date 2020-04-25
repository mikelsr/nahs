package events

import (
	"encoding/json"
	"errors"

	"github.com/mikelsr/bspl"
)

// EventType is used to differentiate events
type EventType string

const (
	// TypeDropEvent an instance was cancelled by one of the
	// parties involved
	TypeDropEvent EventType = "drop"
	// TypeNewEvent a new Instance was created
	TypeNewEvent EventType = "new"
	// TypeUpdateEvent an action was run on an instance
	TypeUpdateEvent EventType = "update"
)

// Event is used by one Node to notify another of a BSPL
// action
type Event interface {
	// Argument of the event. If it is New or Update it
	// will be an Instance, if it is Drop it will be the
	// motive.
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
	Argument    string    `json:"argument"`
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
	case TypeDropEvent, TypeNewEvent, TypeUpdateEvent:
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
	case TypeDropEvent, TypeNewEvent, TypeUpdateEvent:
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
	case TypeDropEvent:
		var a DropEvent
		event, err = a.Unmarshal(marshalledEvent)
		if err != nil {
			return "", err
		}
	case TypeNewEvent:
		var ni NewEvent
		event, err = ni.Unmarshal(marshalledEvent)
		if err != nil {
			return "", err
		}
	case TypeUpdateEvent:
		var nm UpdateEvent
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
	case TypeDropEvent:
		var a DropEvent
		event, err := a.Unmarshal(marshalledEvent)
		if err != nil {
			return err
		}
		a = event.(DropEvent)
		return r.DropInstance(a.InstanceKey(), a.Motive())
	case TypeNewEvent:
		var ni NewEvent
		event, err := ni.Unmarshal(marshalledEvent)
		if err != nil {
			return err
		}
		ni = event.(NewEvent)
		return r.RegisterInstance(ni.Instance())
	case TypeUpdateEvent:
		var nm UpdateEvent
		event, err := nm.Unmarshal(marshalledEvent)
		if err != nil {
			return err
		}
		nm = event.(UpdateEvent)
		return r.UpdateInstance(nm.Instance())
	}
	return nil
}
