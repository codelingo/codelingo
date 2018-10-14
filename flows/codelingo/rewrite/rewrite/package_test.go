package rewrite

import (
	"bytes"

	"testing"

	"github.com/codegangsta/cli"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type cmdSuite struct {
	// jt.FakeHomeSuite
	Context *cli.Context
	stdErr  bytes.Buffer
}

var _ = gc.Suite(&cmdSuite{})
