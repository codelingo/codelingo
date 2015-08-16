package commands

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func (*CMDTest) TestWhichCMD(c *gc.C) {
	cfgPath, closer := testCfg(c)
	defer closer()

	ctx := mockContext(tenetCfgFlg.longArg(), cfgPath, "which")
	c.Assert(WhichCMD.Run(ctx), jc.ErrorIsNil)

	// TODO(waigani) test stdout, regex match: ^/home/[^/]*/\.lingo/tenet\.toml
}
