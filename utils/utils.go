package utils

import (
	"errors"
	"path/filepath"
	"runtime"
)

// GetProjectDir returns the absolute path to the source code of the
// project being run
func GetProjectDir() (string, error) {
	_, fileName, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("Failed to locate project")
	}
	return filepath.Abs(filepath.Dir(fileName))
}
