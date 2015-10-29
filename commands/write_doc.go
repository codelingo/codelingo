package commands

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"

	"github.com/lingo-reviews/lingo/tenet"
)

const defaultTemplate = `# Project Tenets
{{range .All}}
* {{.}}
{{end}}
`

var WriteDocCMD = cli.Command{
	Name:  "write-docs",
	Usage: "output documentation generated from tenets",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "template, t",
			Value:  "",
			Usage:  "path to template file",
			EnvVar: "LINGO_DOC_TEMPLATE",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "tenets.md",
			Usage: "file to write the output to. By default, output file is tenets.md",
		},
	},
	Action: writeDoc,
}

func writeDoc(c *cli.Context) {
	// Find every applicable tenet for this project
	cfgPath, err := tenetCfgPath(c)
	if err != nil {
		oserrf(err.Error())
		return
	}
	cfg, err := buildConfig(cfgPath, CascadeBoth)
	if err != nil {
		oserrf(err.Error())
		return
	}

	r := strings.NewReplacer("/", "_")

	var ts = make(map[string]tenet.Tenet)
	for _, tenetData := range cfg.Tenets {
		// Try to get any installed tenet with matching name
		t, err := tenet.Any(c, tenetData.Name)
		if err != nil {
			// Otherwise try the driver specified in config
			t, err := tenet.New(c, tenetData)
			if err != nil {
				oserrf(err.Error())
				return
			}
			if err = t.InitDriver(); err != nil {
				oserrf(err.Error())
				return
			}
			if err = t.Pull(); err != nil {
				oserrf(err.Error())
				return
			}
		}
		ts[r.Replace(tenetData.Name)] = t
	}

	file, err := os.Create(c.String("output"))
	if err != nil {
		oserrf(err.Error())
		return
	}

	// Add the description of every tenet to the var map and special All array
	// TODO: Add keys for each tenet group
	v := make(map[string]interface{})
	v["All"] = []string{}
	for _, t := range ts {
		n := r.Replace(t.String())
		d, err := t.Description()
		if err != nil {
			oserrf(err.Error())
			return
		}

		v["All"] = append(v["All"].([]string), d)
		v[n] = d
	}

	src := c.String("template")
	if src == "" {
		src = cfg.Template
	}

	var tpl *template.Template
	if src != "" {
		tpl, err = template.ParseFiles(src)
	} else {
		tpl, err = template.New("doc template").Parse(defaultTemplate)
	}
	if err != nil {
		oserrf(err.Error())
		return
	}
	// TODO: Uncomment this when available - need go version update?
	// tpl.Option("missingkey=error") // Error on missing tenets

	if err = tpl.Execute(file, v); err != nil {
		oserrf(err.Error())
	}

	fmt.Printf("Tenet documentation written to %s\n", c.String("output"))
}
