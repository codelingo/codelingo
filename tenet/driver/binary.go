package driver

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/tenet/driver/binary"
)

// Binary is a tenet driver to execute binary tenets found in ~/.lingo/tenets/<repo>/<tenet>
type Binary struct {
	*Base
}

// Check that a file exists at the expected location and is executable. Setup
// the service, but don't start it.
func (b *Binary) Service() (Service, error) {
	tenetPath := b.binPath()

	file, err := os.Open(tenetPath)
	if err != nil {
		return nil, errors.Trace(err)
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if fi.Mode().Perm()&0x49 == 0 {
		return nil, errors.Errorf("%s not exectuable", tenetPath)
	}

	// Note: the service needs to be started and stopped.
	return binary.NewService(tenetPath), nil
}

func (b *Binary) binPath() string {
	if dir := os.Getenv("LINGO_BIN"); dir != "" {
		return filepath.Join(dir, b.Name)
	}
	return filepath.Join(userHomeDir(), ".lingo", "tenets", b.Name)
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func (b *Binary) EditFilename(filename string) (editedFilename string) {
	var absPath string
	var err error
	if absPath, err = filepath.Abs(filename); err == nil {
		return absPath
	}
	log.Printf("could not get absolute path for %q:%v", filename, err)
	return filename
}
