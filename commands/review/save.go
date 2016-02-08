package review

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/tenet"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
)

// Save writes the output to a file at filePath.
func Save(filePath string, issues []*tenet.Issue) error {

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
		log.Printf("could not make filepaths relative to pwd: %v", err)
	}

	jsonIssues, err := json.Marshal(issues)
	if err != nil {
		return errors.Trace(err)
	}

	// CLI will expand tilde for -output ~/file but not -output=~/file. In the
	// latter case, if we can find the user, expand tilde to their home dir.
	if filePath[:2] == "~/" {
		usr, err := user.Current()
		if err == nil {
			dir := usr.HomeDir + "/"
			filePath = strings.Replace(filePath, "~/", dir, 1)
		}
	}

	err = ioutil.WriteFile(filePath, jsonIssues, os.FileMode(0644))
	if err != nil {
		return errors.Errorf("could not write to file %s: %s", filePath, err.Error())
	}
	return nil
}
