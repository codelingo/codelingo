package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

const cfgName = "tenet.toml"

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
		return a
	}
	toggle = true
	return b
}

var a = `
cascade = true
include = "*"
template = ""

[[tenet_group]]
  name = "default"
  template = ""

  [[tenet_group.tenet]]
    name = "lingoreviews/juju_worker_nostate"
    driver = "binary"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]

  [[tenet_group.tenet]]
    name = "lingoreviews/imports"
    driver = "binary"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]

  [[tenet_group.tenet]]
    name = "lingoreviews/license"
    driver = "binary"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]

  [[tenet_group.tenet]]
    name = "lingoreviews/slasher"
    driver = "binary"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]

`[1:]

var b = `
cascade = true
include = "*"
template = ""

[[tenet_group]]
  name = "default"
  template = ""

  [[tenet_group.tenet]]
    name = "lingoreviews/juju_worker_nostate"
    driver = "binary"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]

  [[tenet_group.tenet]]
    name = "lingoreviews/imports"
    driver = "binary"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]

  [[tenet_group.tenet]]
    name = "lingoreviews/simpleseed"
    driver = "binary"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]

  [[tenet_group.tenet]]
    name = "lingoreviews/unused_args"
    driver = "binary"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]

`[1:]
