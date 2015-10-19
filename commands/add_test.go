package commands

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func (*CMDTest) TestAddCMD(c *gc.C) {
	fName, closer := testCfg(c)
	defer closer()

	newTenet := tenet{Author: "waigani", Name: "test"}
	ctx := mockContext(tenetCfgFlg.longArg(), fName, "add", newTenet.String())

	orig, err := readConfigFile(ctx)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(AddCMD.Run(ctx), jc.ErrorIsNil)

	obtained, err := readConfigFile(ctx)
	c.Assert(err, jc.ErrorIsNil)
	expected := append(orig.Tenets, newTenet)
	c.Assert(obtained.Tenets, gc.DeepEquals, expected)
}

func (*CMDTest) TestAddCMDDoubleFails(c *gc.C) {
	// TODO(waigani) write test
}

func (s *CMDTest) TestAddCMDNoURLFails(c *gc.C) {
	ctx := mockContext("add")

	c.Assert(AddCMD.Run(ctx), jc.ErrorIsNil)
	c.Assert(s.stdErr.String(), gc.Equals, "error: expected 1 argument, got 0")
}

func (*CMDTest) TestAddCMDBadURLFails(c *gc.C) {
	// TODO(waigani) write test
}
