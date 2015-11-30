package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lingo-reviews/tenets/go/dev/tenet"

	"github.com/juju/errors"
)

const cfgName = ".lingo"

var toggle bool

// This will add a config to pwd and every sub dir it's run in.
func main() {
	if err := filepath.Walk(".", walk); err != nil {
		panic(err.Error())
	}
}

func walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		cfgPath := filepath.Join(path, cfgName)
		if err := ioutil.WriteFile("./"+cfgPath, []byte(cfg()), 0664); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func cfg() string {
	if toggle {
		toggle = false
		return a()
	}
	toggle = true
	return b()
}

func a() string {
	return fmt.Sprintf(`
cascade = true
include = "*"
template = ""

[[tenet_group]]
  name = "default"
  template = ""

  [[tenet_group.tenet]]
    name = "lingoreviews/juju_worker_nostate"
    driver = "docker"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]
    	opt="%s"

  [[tenet_group.tenet]]
    name = "lingoreviews/imports"
    driver = "docker"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]
    	opt="%s"

  [[tenet_group.tenet]]
    name = "lingoreviews/license"
    driver = "docker"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]
    	opt="%s"

  [[tenet_group.tenet]]
    name = "lingoreviews/slasher"
    driver = "docker"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]
    	opt="%s"

`[1:],
		tenet.RandString(5),
		tenet.RandString(5),
		tenet.RandString(5),
		tenet.RandString(5),
	)
}

func b() string {
	return fmt.Sprintf(`
cascade = true
include = "*"
template = ""

[[tenet_group]]
  name = "default"
  template = ""

  [[tenet_group.tenet]]
    name = "lingoreviews/juju_worker_nostate"
    driver = "docker"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]
    	opt="%s"

  [[tenet_group.tenet]]
    name = "lingoreviews/imports"
    driver = "docker"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]
    	opt="%s"

  [[tenet_group.tenet]]
    name = "lingoreviews/simpleseed"
    driver = "docker"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]
    	opt="%s"

  [[tenet_group.tenet]]
    name = "lingoreviews/unused_arg"
    driver = "docker"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]
    	opt="%s"

`[1:],
		tenet.RandString(5),
		tenet.RandString(5),
		tenet.RandString(5),
		tenet.RandString(5),
	)
}
