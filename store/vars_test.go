package store

import (
	"path/filepath"

	"github.com/mikelsr/nahs/utils"
)

// test variables
var (
	testPath      = _genTestPath()
	testDBPath    = filepath.Join(testPath, "db", "test.db")
	testBSPLPath  = filepath.Join(testPath, "bspl")
	testBSPLFiles = []string{
		filepath.Join(testBSPLPath, "a.bspl"),
		filepath.Join(testBSPLPath, "x.bspl"),
	}
	testKeysPath = filepath.Join(testPath, "keys")
	testNodeN    = 5
)

func _genTestPath() string {
	dir, err := utils.GetProjectDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(dir, "test")
}
