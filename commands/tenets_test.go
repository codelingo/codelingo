package commands

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func (*CMDTest) TestTenetsCMD(c *gc.C) {
	cfgPath, closer := testCfg(c)
	defer closer()

	ctx := mockContext(tenetCfgFlg.longArg(), cfgPath, "tenets")
	c.Assert(TenetsCMD.Run(ctx), jc.ErrorIsNil)

	// TODO(waigani) test stdout
}
