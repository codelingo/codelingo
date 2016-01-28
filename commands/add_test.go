package commands

import (
	jc "github.com/juju/testing/checkers"
	"github.com/lingo-reviews/lingo/commands/common"
	gc "gopkg.in/check.v1"
)

func (*CMDTest) TestAddCMD(c *gc.C) {
	fName, closer := testCfg(c)
	defer closer()

	newTenet := TenetConfig{
		Name:     "waigani/test",
		Driver:   "docker",
		Registry: "hub.docker.com",
		Options:  make(map[string]interface{}),
	}
	ctx := mockContext(c, common.TenetCfgFlg.LongArg(), fName, "add", newTenet.Name)

	orig, err := readConfigFile(fName)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(AddCMD.Run(ctx), jc.ErrorIsNil)

	obtained, err := readConfigFile(fName)
	c.Assert(err, jc.ErrorIsNil)
	expected := append(orig.AllTenets(), newTenet)
	c.Assert(obtained.AllTenets(), gc.DeepEquals, expected)
}

func (*CMDTest) TestAddCMDDoubleFails(c *gc.C) {
	// TODO(waigani) write test
}

func (s *CMDTest) TestAddCMDNoURLFails(c *gc.C) {
	ctx := mockContext(c, "add")

	c.Assert(AddCMD.Run(ctx), jc.ErrorIsNil)
	c.Assert(s.stdErr.String(), gc.Equals, "error: expected 1 argument, got 0\n")
}

func (*CMDTest) TestAddCMDBadURLFails(c *gc.C) {
	// TODO(waigani) write test
}
