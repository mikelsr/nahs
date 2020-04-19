package net

import (
	"testing"

	log "github.com/ipfs/go-log"
)

func TestMain(m *testing.M) {
	log.SetAllLoggers(log.LevelWarn)
	log.SetLogLevel(LogName, "info")
	m.Run()
}
