package commands

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/lingo-reviews/lingo/commands/common"
)

func (s *cmdSuite) TestRemoveCMD(c *gc.C) {
	cfgPath, closer := common.TestCfg(c)
	defer closer()

	tenetToRemove := common.TenetConfig{Name: "lingo-reviews/license"}
	ctx := common.MockContext(c, common.TenetCfgFlg.LongArg(), cfgPath, "remove", tenetToRemove.Name)

	c.Assert(RemoveCMD.Run(ctx), jc.ErrorIsNil)

	obtained, err := common.ReadConfigFile(cfgPath)
	c.Assert(err, jc.ErrorIsNil)
	for _, t := range obtained.AllTenets() {
		c.Assert(t.Name, gc.Not(gc.Equals), tenetToRemove.Name)
	}
}
