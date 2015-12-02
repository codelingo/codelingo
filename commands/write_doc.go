package commands

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"

	"github.com/codegangsta/cli"

	"github.com/lingo-reviews/tenets/go/dev/api"
)

const defaultTemplate = `# Tenets
{{range .Groups}}
{{.}}{{end}}`

const defaultGroupTemplate = `## {{.GroupName | Title}}
{{range .All}}
* {{.}}
{{end}}
`

var WriteDocCMD = cli.Command{
	Name:  "write-docs",
	Usage: "write documentation generated from tenets to a file",
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

	// Make some text processing functions available
	// TODO: Add more. How can we just give the user all of 'strings'?
	var fnMap = template.FuncMap{"Title": strings.Title}

	if src != "" {
		tpl, err = template.ParseFiles(src)
		// TODO: DEMOWARE - This panics on nil pointer - how to make it work?
		// tpl.Funcs(fnMap)
	} else {
		tpl, err = template.New("doc template").Funcs(fnMap).Parse(fallback)
	}
	if err != nil {
		return nil, err
	}

	// TODO: Uncomment this when available - need go version update?
	// tpl.Option("missingkey=error") // Error on missing tenets

	return tpl, nil
}

// MATT funcs should always return errors. Only use the oserrf in the top cmd funcs.
func writeTenetDoc(c *cli.Context, src string, output string) {
	// Find every applicable tenet for this project
	cfgPath, err := tenetCfgPath(c)
	if err != nil {
		oserrf(err.Error())
		return
	}
	cfg, err := buildConfig(cfgPath, CascadeUp)
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
	defer file.Close()

	r := strings.NewReplacer("/", "_")

	// Add the description of every tenet to the var map and special All array
	// Add keys for each tenet group name
	// DEMOWARE: This structure could be a lot simpler

	var wg sync.WaitGroup
	var ts []TenetMeta
	gs := make(map[string]string)
	for _, group := range cfg.TenetGroups {
		for _, tenetCfg := range group.Tenets {
			wg.Add(1)
			go func(group TenetGroup, tenetCfg TenetConfig) {
				defer wg.Done()

				t, err := newTenet(c, tenetCfg)
				if err != nil {
					oserrf(err.Error())
					return
				}

				s, err := t.OpenService()
				if err != nil {
					oserrf(err.Error())
					return
				}
				defer s.Close()
				info, err := s.Info()
				if err != nil {
					oserrf(err.Error())
					return
				}

				d, err := renderedDescription(info, tenetCfg)
				if err != nil {
					oserrf(err.Error())
					return
				}

				ts = append(ts, TenetMeta{
					VarName:     r.Replace(tenetCfg.Name),
					GroupName:   r.Replace(group.Name),
					Description: d,
				})
				gs[r.Replace(group.Name)] = group.Template
			}(group, tenetCfg)
		}
	}
	wg.Wait()

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

func renderedDescription(info *api.Info, cfg TenetConfig) (string, error) {
	tpl, err := template.New("desc template").Parse(info.Description)
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	if err = tpl.Execute(&rendered, cfg.Options); err != nil {
		return "", err
	}

	return rendered.String(), nil
}
