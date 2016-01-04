package commands

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/codegangsta/cli"

	jt "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type CMDTest struct {
	jt.CleanupSuite
	// jt.FakeHomeSuite
	Context *cli.Context
	stdErr  bytes.Buffer
}

var _ = gc.Suite(&CMDTest{})

// 	c.Context = &cli.Context{
// 			App            *App
// 	Command        Command
// 	flagSet        *flag.FlagSet
// 	globalSet      *flag.FlagSet
// 	setFlags       map[string]bool
// 	globalSetFlags map[string]bool
// 	}
// }

var mockTenetCfg = &config{
	TenetGroups: []TenetGroup{
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

func testCfg(c *gc.C) (cfgPath string, closer func()) {
	f, err := ioutil.TempFile("", "mockTenetCfg")
	c.Assert(err, jc.ErrorIsNil)
	ctx := mockContext(c, tenetCfgFlg.longArg(), f.Name(), "noop")
	c.Assert(writeConfigFile(ctx, mockTenetCfg), jc.ErrorIsNil)
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
func mockContext(c *gc.C, args ...string) *cli.Context {
	set := flag.NewFlagSet("test", 0)
	addGlobalOpts(set)

	c.Assert(set.Parse(args), jc.ErrorIsNil)

	ctx := cli.NewContext(cli.NewApp(), set, nil)
	ctx.Command = cli.Command{Name: ctx.Args().First()}
	return ctx
}

func (s *CMDTest) SetUpSuite(c *gc.C) {
	origExiter := exiter
	exiter = func(code int) {
		//noOp func
	}
	stderr = &s.stdErr

	s.AddSuiteCleanup(func(c *gc.C) {
		exiter = origExiter
		stderr = os.Stderr
	})
}

func (s *CMDTest) SetUpTest(c *gc.C) {
	// cleanout err buffer
	s.stdErr = bytes.Buffer{}
}

func (*CMDTest) TestWriteAndReadTenetCfg(c *gc.C) {
	f, err := ioutil.TempFile("", "mockTenetCfg")
	fName := f.Name()
	defer func() {
		os.Remove(fName)
		f.Close()
	}()
	c.Assert(err, jc.ErrorIsNil)
	ctx := mockContext(c, tenetCfgFlg.longArg(), fName)
	c.Assert(writeConfigFile(ctx, mockTenetCfg), jc.ErrorIsNil)

	obtained, err := readConfigFile(fName)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(obtained, gc.DeepEquals, mockTenetCfg)
}

func (*CMDTest) TestDesiredTenetCfgPath(c *gc.C) {
	ctx := mockContext(c, tenetCfgFlg.longArg(), "custom/cfg/path")
	c.Assert(desiredTenetCfgPath(ctx), gc.Equals, "custom/cfg/path")
}

func (*CMDTest) TestTenetCfgPath(c *gc.C) {
	// TODO(waigani) Do what skip says. init .lingo in tmp dir, set default to it.
	c.Skip("Errors if default .lingo cannot be found.")
	defaultPathRegex := `^/home/.*\.lingo_home/\.lingo`
	for cfgPath, expected := range map[string]string{
		// the following cfg files are not found, so the default is returned.
		"rand/path/no/file/.lingo": defaultPathRegex,
		"./.lingo":                 defaultPathRegex,
		".lingo":                   defaultPathRegex,
		"/.lingo":                  defaultPathRegex,
	} {
		ctx := mockContext(c, tenetCfgFlg.longArg(), cfgPath)
		cPath, err := tenetCfgPath(ctx)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(cPath, gc.Matches, expected)
	}
}

func (*CMDTest) TestMockContext(c *gc.C) {
	ctx := mockContext(c, tenetCfgFlg.longArg(), "custom/cfg/path", "add")
	c.Assert(ctx.Args(), gc.DeepEquals, cli.Args{"add"})
	c.Assert(ctx.GlobalString(tenetCfgFlg.long), gc.Equals, "custom/cfg/path")
	c.Assert(ctx.Command.Name, gc.Equals, "add")
}
