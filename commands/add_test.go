package commands

import (
	jc "github.com/juju/testing/checkers"
	"github.com/lingo-reviews/lingo/commands/common"
	gc "gopkg.in/check.v1"
)

func (*cmdSuite) TestAddCMD(c *gc.C) {
	fName, closer := common.TestCfg(c)
	defer closer()

	newTenet := common.TenetConfig{
		Name:     "waigani/test",
		Driver:   "docker",
		Registry: "hub.docker.com",
		Options:  make(map[string]interface{}),
	}
	ctx := common.MockContext(c, common.TenetCfgFlg.LongArg(), fName, "add", newTenet.Name)

	orig, err := common.ReadConfigFile(fName)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(AddCMD.Run(ctx), jc.ErrorIsNil)

	obtained, err := common.ReadConfigFile(fName)
	c.Assert(err, jc.ErrorIsNil)
	expected := append(orig.AllTenets(), newTenet)
	c.Assert(obtained.AllTenets(), gc.DeepEquals, expected)
}

func (*cmdSuite) TestAddCMDDoubleFails(c *gc.C) {
	// TODO(waigani) write test
}

func (s *cmdSuite) TestAddCMDNoURLFails(c *gc.C) {
	ctx := common.MockContext(c, "add")

	c.Assert(AddCMD.Run(ctx), jc.ErrorIsNil)
	c.Assert(s.stdErr.String(), gc.Equals, "error: expected 1 argument, got 0\n")
}

func (*cmdSuite) TestAddCMDBadURLFails(c *gc.C) {
	// TODO(waigani) write test
}
