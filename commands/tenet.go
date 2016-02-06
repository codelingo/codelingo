package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"

	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/tenet"
	"github.com/lingo-reviews/lingo/util"
	"github.com/lingo-reviews/tenets/go/dev/api"
)

// TenetCMD is a fallthrough CMD which treats command as the tenet name and
// passes through any arguments to the tenet.
func TenetCMD(ctx *cli.Context, command string) {
	var commandIsTenet bool
	var cfg common.TenetConfig
	// Does the command match an installed tenet?
	for _, cfg = range listTenets(ctx) {
		if command == cfg.Name {
			commandIsTenet = true
			break
		}
	}
	if !commandIsTenet {
		fmt.Println("command not found")
		return
	}

	tnCMDs, err := newTenetCMDs(ctx, cfg)
	if err != nil {
		common.OSErrf(err.Error())
		return
	}
	defer tnCMDs.closeService()

	if err := tnCMDs.run(); err != nil {
		common.OSErrf(err.Error())
		return
	}
	return
}

func (c *tenetCMDs) run() error {
	method := "help"
	args := c.ctx.Args()
	if len(args[1:]) > 0 {
		method = args[1]
	}

	switch method {
	case "help":
		args := c.ctx.Args()
		if len(args) > 3 {
			return c.printCmdHelp(args[2])
		}
		return c.printHelp()
	case "info":
		return c.printInfo()
	case "description":

		info, err := c.info()
		if err != nil {
			return errors.Trace(err)
		}
		fmt.Println(info.Description)

	case "review":
		return c.review()
	}

	return errors.Errorf("tenet does not have method %q", method)
}

type tenetCMDs struct {
	ctx     *cli.Context
	cfg     common.TenetConfig
	tn      tenet.Tenet
	service tenet.TenetService
}

func newTenetCMDs(ctx *cli.Context, cfg common.TenetConfig) (*tenetCMDs, error) {
	tn, err := common.NewTenet(cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &tenetCMDs{
		ctx: ctx,
		cfg: cfg,
		tn:  tn,
	}, nil
}

func (c *tenetCMDs) openService() (tenet.TenetService, error) {
	if c.service == nil {
		s, err := c.tn.OpenService()
		if err != nil {
			return nil, errors.Trace(err)
		}
		c.service = s
	}
	return c.service, nil
}

func (c *tenetCMDs) closeService() error {
	if s := c.service; s != nil {
		return s.Close()
	}
	return nil
}

func (c *tenetCMDs) printInfo() error {
	info, err := c.info()
	if err != nil {
		return errors.Trace(err)
	}
	text, err := util.FormatOutput(info, infoTemplate)
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println(text)
	return nil
}

func (c *tenetCMDs) info() (*api.Info, error) {
	s, err := c.openService()
	if err != nil {
		return nil, errors.Trace(err)
	}
	info, err := s.Info()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return info, nil
}

func (c *tenetCMDs) printHelp() error {
	var text string
	var err error

	info, err := c.info()
	if err != nil {
		return errors.Trace(err)
	}
	text, err = fmtHelp(info)
	if err != nil {
		return errors.Trace(err)
	}

	fmt.Println(text)
	return nil
}

func (c *tenetCMDs) review() error {

	return errors.New("not implemented")

	// s, err := c.openService()
	// if err != nil {
	// 	return errors.Trace(err)
	// }

	// s.Review(filesc, issuesc)
}

func fmtHelp(info *api.Info) (string, error) {
	h := struct {
		*api.Info
		Commands []*tenetCommand
	}{
		Info:     info,
		Commands: tenetCommands(),
	}
	return util.FormatOutput(h, helpTemplate)
}

func (c *tenetCMDs) printCmdHelp(cmdName string) error {
	for _, cmd := range tenetCommands() {
		if cmd.Name == cmdName {
			out, err := util.FormatOutput(cmd, cmdHelpTemplate)
			if err != nil {
				return errors.Trace(err)
			}
			fmt.Print(out)
			return nil
		}
	}

	fmt.Printf("no help found for %q", cmdName)
	return nil
}

// TODO(waigani) list commands under help
var helpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   command {{.Name}} [arguments...]{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}
`

var cmdHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   command {{.Name}} [arguments...]{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}
`

var infoTemplate = `NAME:
	{{.Name}}
LANGUAGE:
	{{.Language}}
USAGE:
	{{.Usage}}
VERSION:
	{{.Version}}{{if .Options}}
OPTIONS:
	The following option(s) can be set in .lingo or with the --options flag when
	running a review:
{{range .Options}}
	- {{.Name}} ("{{.Value}}"): {{.Usage}}
{{end}}
{{end}}{{if .Metrics}}
METRICS:
	{{.Metrics}}
{{end}}{{if .Tags}}
TAGS:
	{{.Tags}}
{{end}}
`
