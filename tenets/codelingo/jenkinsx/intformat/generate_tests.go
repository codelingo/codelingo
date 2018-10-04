package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
		var idString string
		var integrationDirective = ""
		if example.hasIntegrationDirective {
			integrationDirective = "// +build integration"
			idString += "d"
		} else {
			idString += "o"
		}

		if example.hasValidPackage {
			idString += "p"
		} else {
			idString += "o"
		}

		if example.hasValidFilename {
			idString += "f"
		} else {
			idString += "o"
		}

		testPackageName := idString
		if example.hasValidPackage {
			testPackageName += "_test"
		} else {
			testPackageName += invalidator
		}

		testFileName := idString
		if example.hasValidFilename {
			testFileName += "_integration_test"
		} else {
			testFileName += invalidator
		}

		testFileName = filepath.Join("src", idString, testFileName+".go")
		dirFileName := filepath.Join("src", idString)
		srcFileName := filepath.Join("src", idString, idString+".go")

		if _, err := os.Stat(dirFileName); os.IsNotExist(err) {
			os.Mkdir(dirFileName, os.FileMode(0777))
		}

		fileMode := os.FileMode(0666)
		err := ioutil.WriteFile(testFileName, []byte(fmt.Sprintf(
			testFileTemplate,
			integrationDirective,
			testPackageName,
			idString,
			idString,
		)), fileMode)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(srcFileName, []byte(fmt.Sprintf(
			srcFileTemplate,
			idString,
			idString,
		)), fileMode)
		if err != nil {
			panic(err)
		}
	}
}

var invalidator = "invalid"

var testFileTemplate = `
%s

package %s

func Test%s() {
	%s()
}
`[1:]

var srcFileTemplate = `
package %s

func %s() {}
`[1:]
