package review

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/juju/errors"
	"github.com/lingo-reviews/tenets/go/dev/api"
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
func Output(outputType OutputFormat, outputPath string, issues []*api.Issue) string {

	// // TODO(waigani) provide a flag to control filepath formatting
	// make paths relative to lingo
	pwd, err := os.Getwd()
	if err == nil {
		for _, i := range issues {
			absStart := i.Position.Start.Filename
			absEnd := i.Position.End.Filename
			prefix := pwd + "/"
			i.Position.Start.Filename = strings.TrimPrefix(absStart, prefix)
			i.Position.End.Filename = strings.TrimPrefix(absEnd, prefix)

			i.Name = i.Position.Start.Filename
		}
	} else {
		log.Printf("could not make output filepaths relative to pwd: %v", err)
	}

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

	err = ioutil.WriteFile(outputPath, b.Bytes(), os.FileMode(0644))
	if err != nil {
		panic(errors.Errorf("could not write to file %s: %s", outputPath, err.Error()))
	}
	return fmt.Sprintf("output written to %s", outputPath)
}

func format(outputFmt OutputFormat, issues []*api.Issue) bytes.Buffer {
	var b bytes.Buffer
	switch outputFmt {
	case plainText:
		b = plainFormat(issues)
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

func plainFormat(issues []*api.Issue) bytes.Buffer {
	var out bytes.Buffer
	digits := len(fmt.Sprintf("%d", len(issues)))
	for n, i := range issues {
		out.WriteString(fmt.Sprintf("%*d. %s:%d\n    %s\n", digits, n+1, i.Name, i.Position.Start.Line, strings.TrimSpace(i.Comment)))
	}
	return out
}
