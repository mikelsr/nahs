package net

import (
	"strings"
)

// ErrHandleEvent is returned when handling
// events and sent as a response
type ErrHandleEvent struct {
	ID     string
	Reason string
}

func (e ErrHandleEvent) Error() string {
	var sb strings.Builder
	sb.WriteString("Could not handle event '" + e.ID + "'")
	if e.Reason != "" {
		sb.WriteString(": " + e.Reason)
	}
	return sb.String()
}
