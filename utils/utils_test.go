package utils

import (
	"regexp"
	"testing"
)

func TestGetProjectDir(t *testing.T) {
	dir, err := GetProjectDir()
	if err != nil {
		t.FailNow()
	}
	match, err := regexp.MatchString("^(/[^/ ]*)+/?$", dir)
	if err != nil || !match {
		t.FailNow()
	}
}
