package flow

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/codelingo/lingo/app/util"
	"github.com/juju/errors"
	"github.com/urfave/cli"
)

// Read a codelingo.yaml file from a filepath argument
func ReadDotLingo(ctx *cli.Context) (string, error) {
	var dotlingo []byte

	if filename := ctx.String(util.LingoFile.Long); filename != "" {
		var err error
		dotlingo, err = ioutil.ReadFile(filename)
		if err != nil {
			return "", errors.Trace(err)
		}
	}
	return string(dotlingo), nil
}

// Repeated code in platform
// TODO: variadic fan in - read only
func ErrFanIn(err1c, err2c chan error) chan error {
	errc := make(chan error)
	go func() {
		for {
			if err1c == nil && err2c == nil {
				break
			}

			select {
			case err, ok := <-err1c:
				if !ok {
					err1c = nil
					continue
				}
				errc <- err
			case err, ok := <-err2c:
				if !ok {
					err2c = nil
					continue
				}
				errc <- err
			}
		}
		close(errc)
	}()
	return errc
}

// TODO(waigani) move this to codelingo/sdk/flow
func HandleErr(err error) {
	if errors.Cause(err).Error() == "ui" {
		if e, ok := err.(*errors.Err); ok {
			log.Println(e.Underlying())
			fmt.Println(e.Underlying())
			os.Exit(1)
		}
	}
	fmt.Println(err.Error())
}

// TODO(waigani) this should live under the VCS domain, not Flows
const NoCommitErrMsg = "This looks like a new repository. Please make an initial commit before running `lingo review`. This is only Required for the initial commit, subsequent changes to your repo will be picked up by lingo without committing."

// TODO(waigani) use typed error
func NoCommitErr(err error) bool {
	return strings.Contains(err.Error(), "ambiguous argument 'HEAD'")
}
