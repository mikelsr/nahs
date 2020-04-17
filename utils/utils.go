package utils

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetProjectDir returns the absolute path to the source code of the
// project being run
func GetProjectDir() (string, error) {
	_, fileName, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("Failed to locate project")
	}
	dir, err := filepath.Abs(filepath.Dir(fileName))
	if err != nil {
		return "", err
	}
	path := strings.Split(dir, string(os.PathSeparator))
	dir = "/" + filepath.Join(path[:len(path)-1]...)
	return dir, nil
}
