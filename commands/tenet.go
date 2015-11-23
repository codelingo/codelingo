package commands

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"text/template"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"

	"github.com/lingo-reviews/dev/api"
	"github.com/lingo-reviews/lingo/tenet"
)

// TenetCMD is a fallthrough CMD which treats command as the tenet name and
// passes through any arguments to the tenet.
func TenetCMD(ctx *cli.Context, command string) {
	var commandIsTenet bool
	var cfg TenetConfig
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
		oserrf(err.Error())
		return
	}
	defer tnCMDs.closeService()

	if err := tnCMDs.run(); err != nil {
		oserrf(err.Error())
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
	cfg     TenetConfig
	tn      tenet.Tenet
	service tenet.TenetService
}

func newTenetCMDs(ctx *cli.Context, cfg TenetConfig) (*tenetCMDs, error) {
	tn, err := newTenet(ctx, cfg)
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
	text, err := formatOutput(info, infoTemplate)
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
	return formatOutput(h, helpTemplate)
}

func (c *tenetCMDs) printCmdHelp(cmdName string) error {
	for _, cmd := range tenetCommands() {
		if cmd.Name == cmdName {
			out, err := formatOutput(cmd, cmdHelpTemplate)
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

func formatOutput(in interface{}, tmplt string) (helpStr string, _ error) {
	out := new(bytes.Buffer)
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(tmplt))
	err := t.Execute(w, in)
	if err != nil {
		return "", err
	}
	w.Flush()

	return out.String(), nil
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
	{{.Version}}
METRICS:
	{{.Metrics}}
TAGS:
	{{.Tags}}
{{if .Options}}OPTIONS:

	The following option(s) can be set in .lingo or with the --options flag when
	running a review:

{{range .Options}}
	{{.Name}} ({{.Value}}) - {{.Usage}}
	{{end}}
{{end}}
`
