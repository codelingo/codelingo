package review

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/juju/errors"
	"github.com/lingo-reviews/dev/tenet"
)

type OutputFormat string

const (
	plainText  OutputFormat = "plain-text"
	jsonPretty OutputFormat = "json-pretty"
	jsonOut    OutputFormat = "json"
)

// func OutputFormat(outputFmt string) *OutputFormat {
// 	var err error
// 	var o OutputFormat
// 	switch outputFmt {
// 	case "plain-text":
// 		o = plainText
// 	case "json-pretty":
// 		o = jsonPretty
// 	case "json":
// 		o = jsonOut
// 	}
// 	return &o
// }

// finalOutput writes the output to a file at outputFile, unless outputFile ==
// "cli", in which case it returns the output.
func Output(outputType OutputFormat, outputPath string, issues []*tenet.Issue) string {
	b := format(outputType, issues)
	if outputPath == "cli" {
		return b.String()
	}

	// CLI will expand tilde for -output ~/file but not -output=~/file. In the
	// latter case, if we can find the user, expand tilde to their home dir.
	if outputPath[:2] == "~/" {
		usr, err := user.Current()
		if err == nil {
			dir := usr.HomeDir + "/"
			outputPath = strings.Replace(outputPath, "~/", dir, 1)
		}
	}

	err := ioutil.WriteFile(outputPath, b.Bytes(), os.FileMode(0775))
	if err != nil {
		panic(errors.Errorf("could not write to file %s: %s", outputPath, err.Error()))
	}
	return fmt.Sprintf("output written to %s", outputPath)
}

func format(outputFmt OutputFormat, issues []*tenet.Issue) bytes.Buffer {
	var b bytes.Buffer
	fileTemplates := map[string]string{}
	switch outputFmt {
	case plainText:
		for _, issue := range issues {
			fileTemplates[issue.Filename()] += "\n" + issue.String()
		}

		if len(fileTemplates) == 0 {
			fmt.Fprintln(&b, "No Issues Found")
		} else {
			fmt.Fprintln(&b, "Issues Found:")
			for _, ft := range fileTemplates {
				fmt.Fprintln(&b, ft)
			}
		}
	case jsonPretty:
		formatted, err := json.MarshalIndent(issues, "", "\t")
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(&b, string(formatted))
	case jsonOut:
		formatted, err := json.Marshal(issues)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(&b, string(formatted))
	default:
		panic(errors.Errorf("Unrecognised output format %q", outputFmt))
	}
	return b
}
