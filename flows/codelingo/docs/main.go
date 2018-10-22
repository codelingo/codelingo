package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/codelingo/codelingo/flows/codelingo/docs/docs"
	"github.com/codelingo/codelingo/flows/codelingo/docs/docs/parse"
	"github.com/codelingo/codelingo/flows/codelingo/docs/render"
	flowutil "github.com/codelingo/codelingo/sdk/flow"
	"github.com/codelingo/lingo/app/util"
	"github.com/juju/errors"
	"github.com/urfave/cli"
	//	"gopkg.in/russross/blackfriday.v2"
	"github.com/russross/blackfriday"
)

var docsApp = &flowutil.CLIApp{
	App: cli.App{
		Name:    "docs",
		Usage:   "Generate documentation from Tenets",
		Version: "0.0.0",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "output, o",
				Usage: "File to write output to. If no file is given, docs are printed to the terminal.",
			},
			cli.StringFlag{
				Name:  "template, t",
				Value: "default",
				Usage: "The template file to use when generating docs.",
			},
			cli.StringFlag{
				Name:  "format, f",
				Usage: "Format to render docs to. Options are: md and html",
			},
			cli.BoolFlag{
				Name:  "web, w",
				Usage: "Render documentation as a static website and launch it",
			},
		},
		Action: docsAction,
	},
}

func main() {

	fRunner := flowutil.NewFlow(docsApp, nil)
	_, err := fRunner.Run()
	if err != nil {
		util.Logger.Debugw("Review Flow", "err_stack", errors.ErrorStack(err))
		util.FatalOSErr(err)
		return
	}
}

func docsAction(ctx *cli.Context) {
	docsStr, err := formattedDocsFromTenets(ctx)
	if err != nil {

		// Debugging
		util.Logger.Debugw("docsAction", "err_stack", errors.ErrorStack(err))
		util.FatalOSErr(err)
		return
	}

	if err := renderDocs(ctx, docsStr); err != nil {

		// Debugging
		util.Logger.Debugw("docsAction", "err_stack", errors.ErrorStack(err))
		util.FatalOSErr(err)
		return
	}

}

func renderDocs(ctx *cli.Context, docStr string) error {

	if ctx.Bool("web") {
		return errors.New("Not Implemented")
	}

	if tgtFile := ctx.String("output"); tgtFile != "" {
		if err := ioutil.WriteFile(tgtFile, []byte(docStr), 0644); err != nil {
			return err
		}

		fmt.Println("Success! Docs written to " + tgtFile)
		return nil
	}

	_, err := os.Stdout.Write([]byte(docStr))
	return errors.Trace(err)
}

func formattedDocsFromTenets(ctx *cli.Context) (string, error) {

	allDocs, err := docMaps(ctx)
	if err != nil {
		return "", errors.Trace(err)
	}

	var templateStr []byte
	if templateFile := ctx.String("template"); templateFile != "" {
		templateStr, err = ioutil.ReadFile(templateFile)
	}

	mdDocs, err := tenetsToRawMdDocs(allDocs, string(templateStr))
	if err != nil {
		return "", errors.Trace(err)
	}

	docFmt := ctx.String("format")
	switch docFmt {
	case "":

		// If no format is specified, default to the format for the render medium

		if tgtFile := ctx.String("output"); tgtFile != "" {
			// if writing to file, default to md
			return mdDocs, nil
		}

		if ctx.Bool("web") {
			// if webpage, default to html
			return string(blackfriday.MarkdownCommon([]byte(mdDocs))), nil
		}

		// otherwise, we're catting to the terminal
		renderer, extensions := render.Terminal()
		return string(blackfriday.Markdown([]byte(mdDocs), renderer, extensions)), nil

	case "md":
		return mdDocs, nil
	case "html":
		return string(blackfriday.MarkdownCommon([]byte(mdDocs))), nil
	}
	return "", errors.New("unknown format: " + docFmt)
}

func docMaps(ctx *cli.Context) ([]map[string]string, error) {

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, errors.Trace(err)
	}

	files, err := docs.GetLingoFiles(workingDir)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var allDocs []map[string]string
	for _, fileSRC := range files {

		lingoFiles, err := parse.Parse(string(fileSRC))
		if err != nil {
			return nil, errors.Trace(err)
		}

		for _, lingoFile := range lingoFiles {
			for _, tenet := range lingoFile.Tenets {

				for _, flow := range tenet.Bots {
					if flow.Owner == "codelingo" && flow.Name == "docs" {
						allDocs = append(allDocs, flow.Config)
					}
				}

			}
		}
	}

	return allDocs, nil
}

// TODO(waigani) in cmd help, show default template
var defaultTemplate string = `
# Contributor Guide
{{range .}}
## {{.title}}

{{.body}}
{{end}}
`[1:]

func tenetsToRawMdDocs(docs []map[string]string, templateSRC string) (string, error) {
	tempSRC := defaultTemplate
	if templateSRC != "" {
		tempSRC = templateSRC
	}
	t := template.Must(template.New("t3").Parse(tempSRC))

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, docs); err != nil {
		return "", errors.Trace(err)
	}

	return tpl.String(), nil
}
