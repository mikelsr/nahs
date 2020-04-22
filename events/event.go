package events

import (
	"encoding/json"

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
	Argument() interface{}
	ID() string
	Instance() bspl.Instance
	Type() EventType

	Marshal() []byte
	UnMarshal([]byte) (Event, error)
}

// EventWrapper is used by different event types
// to marshal themselves
type EventWrapper struct {
	Argument string    `json:"arguments"`
	ID       string    `json:"id"`
	Instance string    `json:"instance"`
	Type     EventType `json:"event_type"`
}

// Marshal an EventWrapper
func (e EventWrapper) Marshal() ([]byte, error) {
	return json.Marshal(e)
}
