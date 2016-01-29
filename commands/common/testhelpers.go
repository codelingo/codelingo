package common

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

// stderr is a var for mocking in tests
var Stderr io.Writer = os.Stderr

// exiter is a var for mocking in tests
var Exiter = func(code int) {
	os.Exit(code)
}

func TestCfg(c *gc.C) (cfgPath string, closer func()) {
	f, err := ioutil.TempFile("", "MockTenetCfg()")
	c.Assert(err, jc.ErrorIsNil)
	ctx := MockContext(c, TenetCfgFlg.LongArg(), f.Name(), "noop")
	c.Assert(WriteConfigFile(ctx, MockTenetCfg()), jc.ErrorIsNil)
	return f.Name(), func() {
		os.Remove(f.Name())
		f.Close()
	}
}

func longName(f cli.Flag) string {
	parts := strings.Split(f.String(), ",")
	return strings.TrimLeft(parts[0], "-")
}

func addGlobalOpts(set *flag.FlagSet) {
	for _, flg := range GlobalOptions {
		lName := longName(flg)
		switch f := flg.(type) {
		case cli.BoolFlag:
			set.Bool(lName, false, f.Usage)
		case cli.StringFlag:
			set.String(lName, f.Value, f.Usage)
		}
	}
}

// mockContext is a test helper for testing commands. Flags should only be set
// with their long name.
func MockContext(c *gc.C, args ...string) *cli.Context {
	set := flag.NewFlagSet("test", 0)
	addGlobalOpts(set)

	c.Assert(set.Parse(args), jc.ErrorIsNil)

	ctx := cli.NewContext(cli.NewApp(), set, nil)
	ctx.Command = cli.Command{Name: ctx.Args().First()}
	return ctx
}

func MockTenetCfg() *Config {
	return &Config{TenetGroups: []TenetGroup{
		{Name: "default",
			Tenets: []TenetConfig{
				{
					Name: "lingoreviews/tenetseed:latest",
				}, {
					Name: "lingoreviews/space_after_forward_slash",
				}, {
					Name: "lingo-reviews/unused_function_args",
				}, {
					Name: "lingo-reviews/license",
					Options: map[string]interface{}{
						"header": "// MIT\n",
					},
				},
			},
		},
	},
	}
}
