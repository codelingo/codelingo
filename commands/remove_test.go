package commands

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func (s *CMDTest) TestRemoveCMD(c *gc.C) {
	cfgPath, closer := testCfg(c)
	defer closer()

	tenetToRemove := TenetConfig{Name: "lingo-reviews/license"}
	ctx := mockContext(c, tenetCfgFlg.longArg(), cfgPath, "remove", tenetToRemove.Name)

	c.Assert(RemoveCMD.Run(ctx), jc.ErrorIsNil)

	obtained, err := readConfigFile(cfgPath)
	c.Assert(err, jc.ErrorIsNil)
	for _, t := range obtained.AllTenets() {
		c.Assert(t.Name, gc.Not(gc.Equals), tenetToRemove.Name)
	}
}
