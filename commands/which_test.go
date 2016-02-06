package commands

import (
	jc "github.com/juju/testing/checkers"
	"github.com/lingo-reviews/lingo/commands/common"
	gc "gopkg.in/check.v1"
)

func (*cmdSuite) TestWhichCMD(c *gc.C) {
	cfgPath, closer := common.TestCfg(c)
	defer closer()

	ctx := common.MockContext(c, common.TenetCfgFlg.LongArg(), cfgPath, "which")
	c.Assert(WhichCMD.Run(ctx), jc.ErrorIsNil)

	// TODO(waigani) test stdout, regex match: ^/home/[^/]*/\.lingo/tenet\.toml
}
