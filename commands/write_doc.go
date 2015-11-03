package commands

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"

	"github.com/lingo-reviews/lingo/tenet"
)

const defaultTemplate = `# Tenets
{{range .Groups}}
{{.}}{{end}}`

const defaultGroupTemplate = `## {{.GroupName}}
{{range .All}}
* {{.}}
{{end}}
`

var WriteDocCMD = cli.Command{
	Name:  "write-docs",
	Usage: "write output documentation generated from tenets to a file",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "template, t",
			Value:  "",
			Usage:  "path to template file",
			EnvVar: "LINGO_DOC_TEMPLATE",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "lingo_docs/tenets.md",
			Usage: "file to write the output to. By default, output file is lingo_docs/tenets.md",
		},
	},
	Action: writeDoc,
}

type TenetMeta struct {
	VarName     string
	GroupName   string
	Description string
}

func writeDoc(c *cli.Context) {
	output := c.String("output")
	writeTenetDoc(c, c.String("template"), output)
	fmt.Printf("Tenet documentation written to %s\n", output)
}

func makeTemplate(src string, fallback string) (*template.Template, error) {
	var tpl *template.Template
	var err error
	if src != "" {
		tpl, err = template.ParseFiles(src)
	} else {
		tpl, err = template.New("doc template").Parse(fallback)
	}
	if err != nil {
		return nil, err
	}

	// TODO: Uncomment this when available - need go version update?
	// tpl.Option("missingkey=error") // Error on missing tenets

	return tpl, nil
}

func writeTenetDoc(c *cli.Context, src string, output string) {
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

	if dir, _ := path.Split(output); dir != "" {
		err := os.MkdirAll(dir, 0775)
		if err != nil {
			oserrf(err.Error())
			return
		}
	}

	file, err := os.Create(output)
	if err != nil {
		oserrf(err.Error())
		return
	}

	r := strings.NewReplacer("/", "_")

	// Add the description of every tenet to the var map and special All array
	// Add keys for each tenet group name
	// DEMOWARE: This structure could be a lot simpler
	var ts []TenetMeta
	gs := make(map[string]string)
	for _, group := range cfg.TenetGroups {
		for _, tenetData := range group.Tenets {
			// Try to get any installed tenet with matching name
			t, err := tenet.Any(c, tenetData.Name, tenetData.Options)
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
			d, err := tenet.RenderedDescription(t)
			if err != nil {
				oserrf(err.Error())
				return
			}
			ts = append(ts, TenetMeta{
				VarName:     r.Replace(tenetData.Name),
				GroupName:   r.Replace(group.Name),
				Description: d,
			})
			gs[r.Replace(group.Name)] = group.Template
		}
	}

	// Make the description available in multiple places:
	// Top level template
	// ├── All (array)
	// │   ├── desc1
	// │   └── desc2
	// ├── Groups (array)
	// │   ├── renderedGroup1
	// │   └── renderedGroup2
	// ├── tenetName1 (string) desc1
	// ├── tenetName2 (string) desc2
	// ├── groupName1 (string) renderedGroup1
	// └── groupName2 (string) renderedGroup2
	//
	// Group template
	// ├── All (array)
	// │   └── desc1
	// ├── GroupName (string) name from config
	// └── tenetName1 (string) desc1

	// Render groups first
	renderedGroups := make(map[string]string)
	for n, tmpl := range gs {
		g := make(map[string]interface{})
		g["All"] = []string{}
		g["GroupName"] = n
		for _, tm := range ts {
			if tm.GroupName == n {
				g["All"] = append(g["All"].([]string), tm.Description)
				g[tm.VarName] = tm.Description
			}
		}

		tpl, err := makeTemplate(tmpl, defaultGroupTemplate)
		if err != nil {
			oserrf(err.Error())
			return
		}

		var rg bytes.Buffer
		if err = tpl.Execute(&rg, g); err != nil {
			oserrf(err.Error())
			return
		}
		renderedGroups[n] = rg.String()
	}

	v := make(map[string]interface{})
	v["All"] = []string{}
	v["Groups"] = renderedGroups
	for _, tm := range ts {
		v["All"] = append(v["All"].([]string), tm.Description)
		v[tm.VarName] = tm.Description
		for n := range gs {
			v[n] = renderedGroups[n]
		}
	}

	if src == "" {
		src = cfg.Template
	}
	tpl, err := makeTemplate(src, defaultTemplate)
	if err != nil {
		oserrf(err.Error())
		return
	}

	if err = tpl.Execute(file, v); err != nil {
		oserrf(err.Error())
	}
}
