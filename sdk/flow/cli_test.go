package flow

import (
	"github.com/golang/protobuf/proto"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/urfave/cli"
)

func (s *flowSuite) TestDecoratorInput(c *gc.C) {

	decoratorCMD.ConfirmDecorated = func(ctx *cli.Context, msg proto.Message) (bool, error) {

		c.Assert(ctx.Bool("bool-flag"), gc.Equals, true)
		c.Assert(ctx.String("some-decorator-flag"), gc.Equals, "arg1")
		c.Assert(ctx.Args(), jc.DeepEquals, cli.Args{"\"lit1\"", "arg2", "\"lit2\""})

		return true, nil
	}

	input := "--bool-flag --some-decorator-flag arg1 \"lit1\" arg2 \"lit2\""
	// ctx, err := NewCtx(cliCMD.Command, strings.Split(input, " ")[1:])
	// c.Assert(err, jc.ErrorIsNil)

	fRunner := NewFlow(cliCMD, decoratorCMD)

	_, err := fRunner.ConfirmDecorated(input, nil)
	c.Assert(err, jc.ErrorIsNil)
}

var cliCMD = &CLIApp{
	App: cli.App{
		Name:  "cli",
		Usage: "Dummy cli cmd",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "some-cli-flag",
				Usage: "some cli usage",
			},
		},
	},
	Request: func(*cli.Context) (chan proto.Message, <-chan *UserVar, chan error, func(), error) {
		return nil, nil, nil, nil, nil
	},
}

var decoratorCMD = &DecoratorApp{
	App: cli.App{
		Name:  "decorator",
		Usage: "Dummy decorator cmd",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "some-decorator-flag",
				Usage: "some decorator usage",
			},
			cli.StringFlag{
				Name:  "not-used-flag",
				Usage: "some decorator usage",
			},
			cli.BoolFlag{
				Name:  "bool-flag",
				Usage: "some decorator usage",
			},
		},
	},
}

func (s *flowSuite) TestDecoratorArgs(c *gc.C) {
	for _, test := range decoratorArgsTest {
		result := DecoratorArgs(test.input)
		c.Assert(result, jc.DeepEquals, test.expected)
	}
}

var decoratorArgsTest = []struct {
	input    string
	expected []string
}{
	{
		input:    "review -f -someflag",
		expected: []string{"-f", "-someflag"},
	},
	{
		input:    "review -f \"{{something}}\"",
		expected: []string{"-f", "\"{{something}}\""},
	},
	{
		input:    `rewrite -r "errors.Errorf(\"{{formatString}})\""`,
		expected: []string{"-r", `"errors.Errorf(\"{{formatString}})\""`},
	},
}
