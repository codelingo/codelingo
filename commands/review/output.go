package review

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"text/template"

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
	switch outputFmt {
	case plainText:
		if len(issues) == 0 {
			fmt.Fprintln(&b, "No issues found")
			break
		}
		for _, issue := range issues {
			fmt.Fprintln(&b, FormatPlainText(issue))
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

// Comment returns the comment for the issue, depending on the context of the issue.
func Comment(issue *tenet.Issue, commSet *tenet.CommentSet) (string, error) {
	comments := commSet.CommentsForContext(issue.Context)
	if len(comments) == 0 {
		comments = commSet.CommentsForContext(tenet.All)
	}

	// build comments with template args
	t := template.New("comment template")
	// default message if no comment set
	commentTemplate := "Issue Found"
	if len(comments) > 0 {
		commentTemplate = comments[0].Template
	}
	ct, err := t.Parse(commentTemplate) // TODO(waigani) This only returns the first comment for each context.
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = ct.Execute(&b, issue.CommVars)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
