package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/juju/errors"
)

const cfgName = "tenet.toml"

// This will add a config to pwd and every sub dir it's run in.
func main() {
	if err := filepath.Walk(".", walk); err != nil {
		panic(err.Error())
	}
}

func walk(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		cfgPath := filepath.Join(path, cfgName)
		if err := ioutil.WriteFile("./"+cfgPath, []byte(s), 0664); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

var s = `
cascade = true
include = "*"
template = ""

[[tenet_group]]
  name = "default"
  template = ""

  [[tenet_group.tenet]]
    name = "lingoreviews/simpleseed"
    driver = "binary"
    registry = ""
    tag = ""
    [tenet_group.tenet.options]
`[1:]
