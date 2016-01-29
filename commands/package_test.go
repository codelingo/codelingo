package commands

import (
	"bytes"
	"os"

	"testing"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/commands/common"

	jt "github.com/juju/testing"
	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type cmdSuite struct {
	jt.CleanupSuite
	// jt.FakeHomeSuite
	Context *cli.Context
	stdErr  bytes.Buffer
}

var _ = gc.Suite(&cmdSuite{})

func (s *cmdSuite) SetUpSuite(c *gc.C) {
	origExiter := common.Exiter
	common.Exiter = func(code int) {
		//noOp func
	}
	common.Stderr = &s.stdErr

	s.AddSuiteCleanup(func(c *gc.C) {
		common.Exiter = origExiter
		common.Stderr = os.Stderr
	})
}

func (s *cmdSuite) SetUpTest(c *gc.C) {
	// cleanout err buffer
	s.stdErr = bytes.Buffer{}
}
