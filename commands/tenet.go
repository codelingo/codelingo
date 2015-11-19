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
	"github.com/lingo-reviews/lingo/tenet/driver"
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

	if err := runTenetCMD(ctx, command, cfg); err != nil {
		oserrf(err.Error())
	}
	return
}

func runTenetCMD(ctx *cli.Context, command string, cfg TenetConfig) error {
	var method string
	args := ctx.Args()
	if len(args[1:]) > 0 {
		method = args[1]
	}

	switch method {
	case "help", "":
		var text string
		var err error
		if len(args) > 3 {
			text, err = cmdHelp(args[2])
		} else {
			info, err := tenetInfo(ctx, cfg)
			if err != nil {
				return errors.Trace(err)
			}
			text, err = help(info)
		}
		if err != nil {
			return errors.Annotatef(err, "error running method %q", method)
		}

		fmt.Println(text)
	case "info":
		info, err := tenetInfo(ctx, cfg)
		if err != nil {
			return errors.Trace(err)
		}

		text, err := formatOutput(info, infoTemplate)
		if err != nil {
			return errors.Trace(err)
		}

		fmt.Println(text)
	case "description":

		info, err := tenetInfo(ctx, cfg)
		if err != nil {
			return errors.Trace(err)
		}

		fmt.Println(info.Description)
	case "review":

		// s, err := openService(ctx, cfg)
		// if err != nil {
		// 	return errors.Trace(err)
		// }
		// defer s.Stop()
		// s.Review(filesc, issuesc)
		fmt.Println("not implemented")
	default:
		return errors.Errorf("tenet does not have method %q", method)
	}
	return nil
}

// openService returns a started tenet micro-service. It needs to be
// explicitly stopped.
func openService(ctx *cli.Context, tenetCfg TenetConfig) (tenet.TenetService, error) {
	t, err := tenet.New(ctx, &driver.Base{
		Name:          tenetCfg.Name,
		Driver:        tenetCfg.Driver,
		Registry:      tenetCfg.Registry,
		Tag:           tenetCfg.Tag,
		ConfigOptions: tenetCfg.Options,
	})
	if err != nil {
		return nil, errors.Trace(err)
	}
	s, err := t.Service()
	if err != nil {
		return nil, errors.Trace(err)
	}
	err = s.Start()
	return s, err
}

func tenetInfo(ctx *cli.Context, cfg TenetConfig) (*api.Info, error) {
	s, err := openService(ctx, cfg)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer s.Stop()
	return s.Info()
}

func help(info *api.Info) (string, error) {
	h := struct {
		*api.Info
		Commands []*tenetCommand
	}{
		Info:     info,
		Commands: tenetCommands(),
	}
	return formatOutput(h, helpTemplate)
}

func cmdHelp(cmdName string) (string, error) {
	for _, cmd := range tenetCommands() {
		if cmd.Name == cmdName {
			return formatOutput(cmd, cmdHelpTemplate)
		}
	}
	return fmt.Sprintf("no help found for %q", cmdName), nil
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
