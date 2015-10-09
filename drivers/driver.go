package drivers

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"

	devTenet "github.com/lingo-reviews/dev/tenet"
	"github.com/lingo-reviews/lingo/tenet"
)

func decodeResult(name string, result string) (*tenet.ReviewResult, error) {
	reviewResult := &tenet.ReviewResult{}
	err := json.Unmarshal([]byte(result), reviewResult)
	reviewResult.TenetName = name
	return reviewResult, err
}

type Driver interface {
	Review(t *tenet.Tenet, args ...string) (*tenet.ReviewResult, error)
	Help(t *tenet.Tenet, args ...string) (string, error)
	Version(t *tenet.Tenet) (string, error)
	CommentSet(t *tenet.Tenet) (*devTenet.CommentSet, error)
	Debug(t *tenet.Tenet, args ...string) string
}

func New(d string, c *cli.Context) (Driver, error) {
	switch d {
	case "binary":
		return newBinary(c)
	case "docker":
		return newDocker(c)
	}
	return nil, fmt.Errorf("Invalid driver specified: %q", d)
}
