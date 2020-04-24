package storage

import (
	"testing"
)

func TestMain(m *testing.M) {
	generateDB()
	//generatePeers()
	m.Run()
}
