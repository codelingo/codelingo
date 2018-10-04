package main

import (
	"fmt"
	"io/ioutil"
)

type srccode struct {
	hasIntegrationDirective bool
	hasValidFilename        bool
	hasValidPackage         bool
}

// Integration tests must have a build directive, a certain filename, and a certain
// package name. Generate a test file for each possible combination of having/missing
//  those characteristics.
func main() {
	for _, example := range []srccode{
		{false, false, false},
		{false, false, true},
		{false, true, false},
		{false, true, true},
		{true, false, false},
		{true, false, true},
		{true, true, false},
		{true, true, true},
	} {
		var fname string
		var intd = ""
		if example.hasIntegrationDirective {
			intd = "// +build integration"
			fname += "D"
		} else {
			fname += "O"
		}

		var packageName = ""
		if example.hasValidPackage {
			packageName = "_integration_test"
			fname += "P"
		} else {
			fname += "O"
		}

		if example.hasValidFilename {
			fname += "F"
			fname += "_integration_test"
		} else {
			fname += "O"
		}

		err := ioutil.WriteFile(fname+".go", []byte(fmt.Sprintf(filetemplate, intd, packageName, fname)), 0644)
		if err != nil {
			panic(err)
		}
	}
}

var filetemplate = `
%s

package intformat%s

func %s() {}
`[1:]
