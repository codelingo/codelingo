package flow

import (
	"bytes"

	"testing"

	"github.com/urfave/cli"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type flowSuite struct {
	// jt.FakeHomeSuite
	Context *cli.Context
	stdErr  bytes.Buffer
}

var _ = gc.Suite(&flowSuite{})
