package commands

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func (*CMDTest) TestUpdateCMD(c *gc.C) {
	cfgPath, closer := testCfg(c)
	defer closer()

	ctx := mockContext(c, tenetCfgFlg.longArg(), cfgPath, "update")
	c.Assert(WhichCMD.Run(ctx), jc.ErrorIsNil)

	// TODO(waigani) tenetHome with fake dir
}
