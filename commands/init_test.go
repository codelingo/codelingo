package commands

import (
	"flag"

	"github.com/codegangsta/cli"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func (s *CMDTest) TestInitCMD(c *gc.C) {

	app := cli.NewApp()
	set := flag.NewFlagSet("test", 0)
	test := []string{"pwd"}
	set.Parse(test)

	ctx := cli.NewContext(app, set, nil)
	// ctx.GlobalString("name")

	c.Assert(InitCMD.Run(ctx), jc.ErrorIsNil)

}
