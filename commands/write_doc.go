package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
	outputPath := c.String("output")
	if dir, _ := path.Split(outputPath); dir != "" {
		err := os.MkdirAll(dir, 0775)
		if err != nil {
			oserrf(err.Error())
			return
		}
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		oserrf(err.Error())
		return
	}
	defer outputFile.Close()

	if err := writeTenetDoc(c, c.String("template"), outputFile); err != nil {
		oserrf(err.Error())
		return
	}
	fmt.Printf("Tenet documentation written to %s\n", outputFile.Name())
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

func writeTenetDoc(c *cli.Context, src string, w io.Writer) error {
	// Find every applicable tenet for this project
	cfgPath, err := tenetCfgPath(c)
	if err != nil {
		return err
	}
	cfg, err := buildConfig(cfgPath, CascadeUp)
	if err != nil {
		return err
	}

	r := strings.NewReplacer("/", "_")

	// Add the description of every tenet to the var map and special All array
	// Add keys for each tenet group name
	// DEMOWARE: This structure could be a lot simpler

	type result struct {
		name     string
		template string
		err      error
	}

	var wg sync.WaitGroup
	var ts []TenetMeta
	var results []result
	for _, group := range cfg.TenetGroups {
		for range group.Tenets {
			results = append(results, result{
				name: r.Replace(group.Name),
			})
		}
	}
	var i int
	for _, group := range cfg.TenetGroups {
		for _, tenetCfg := range group.Tenets {
			wg.Add(1)
			go func(group TenetGroup, tenetCfg TenetConfig, result *result) {
				defer wg.Done()

				t, err := newTenet(c, tenetCfg)
				if err != nil {
					result.err = err
					return
				}

				s, err := t.OpenService()
				if err != nil {
					result.err = err
					return
				}
				defer s.Close()
				info, err := s.Info()
				if err != nil {
					result.err = err
					return
				}

				d, err := renderedDescription(info, tenetCfg)
				if err != nil {
					result.err = err
					return
				}

				ts = append(ts, TenetMeta{
					VarName:     r.Replace(tenetCfg.Name),
					GroupName:   r.Replace(group.Name),
					Description: d,
				})
				result.template = group.Template
			}(group, tenetCfg, &results[i])
			i++
		}
	}
	wg.Wait()

	// If any of the tenets could not be rendered, return an error now.
	var errs []error
	for _, result := range results {
		if result.err != nil {
			errs = append(errs, result.err)
		}
	}
	switch len(errs) {
	case 0:
	case 1:
		return errs[0]
	default:
		errorStrings := make([]string, len(errs))
		for i, err := range errs {
			errorStrings[i] = err.Error()
		}
		return errors.New(strings.Join(errorStrings, "\n"))
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
	for _, result := range results {
		g := make(map[string]interface{})
		g["All"] = []string{}
		g["GroupName"] = result.name
		for _, tm := range ts {
			if tm.GroupName == result.name {
				g["All"] = append(g["All"].([]string), tm.Description)
				g[tm.VarName] = tm.Description
			}
		}

		tpl, err := makeTemplate(result.template, defaultGroupTemplate)
		if err != nil {
			return err
		}

		var rg bytes.Buffer
		if err = tpl.Execute(&rg, g); err != nil {
			return err
		}
		renderedGroups[result.name] = rg.String()
	}

	v := make(map[string]interface{})
	v["All"] = []string{}
	v["Groups"] = renderedGroups
	for _, tm := range ts {
		v["All"] = append(v["All"].([]string), tm.Description)
		v[tm.VarName] = tm.Description
		for _, result := range results {
			v[result.name] = renderedGroups[result.name]
		}
	}

	if src == "" {
		src = cfg.Template
	}
	tpl, err := makeTemplate(src, defaultTemplate)
	if err != nil {
		return err
	}

	return tpl.Execute(w, v)
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
