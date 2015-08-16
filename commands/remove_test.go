package commands

import (
	jc "github.com/juju/testing/checkers"
	"github.com/waigani/xxx"
	gc "gopkg.in/check.v1"
)

func (*CMDTest) TestRemoveCMD(c *gc.C) {
	cfgPath, closer := testCfg(c)
	defer closer()
	xxx.Print(cfgPath)
	tenetToRemove := tenet{Author: "lingo-reviews", Name: "license"}
	ctx := mockContext(tenetCfgFlg.longArg(), cfgPath, "remove", tenetToRemove.String())

	c.Assert(RemoveCMD.Run(ctx), jc.ErrorIsNil)

	obtained, err := readTenetCfgFile(ctx)
	c.Assert(err, jc.ErrorIsNil)
	for _, t := range obtained.Tenets {
		c.Assert(t.String(), gc.Not(gc.Equals), tenetToRemove.String())
	}
}
