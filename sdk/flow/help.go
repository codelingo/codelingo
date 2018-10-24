package flow

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/urfave/cli"
)

func printHelp(out io.Writer, templ string, data interface{}) {
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	w := tabwriter.NewWriter(out, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		// If the writer is closed, t.Execute will fail, and there's nothing
		// we can do to recover.
		if os.Getenv("CLI_TEMPLATE_ERROR_DEBUG") != "" {
			fmt.Fprintf(cli.ErrWriter, "CLI TEMPLATE ERROR: %#v\n", err)
		}
		return
	}
	w.Flush()

}

var CLIAPPHELPTMP = `

 --- {{.Tagline}} - https://codelingo.io/flow/codelingo/{{.Name}} ---


  FLOW HELP
  =========

   {{.Usage}}


   USAGE
   -----
 
    $ lingo run {{.Name}} [options] 
 

   OPTIONS
   -------
    {{range .VisibleFlags}}
    {{.}}
    {{end}}
`[1:]

var DECAPPHELPTMP = `

  DECORATOR HELP
  ==============

   USAGE
   -----
  
    @{{.Name}} {{.DecoratorUsage}}
    
   {{if len .VisibleFlags}}
   OPTIONS
   -------
    {{range .VisibleFlags}}
    {{.}}
    {{end}}
   {{end}}
   EXAMPLE
   -------
 
    # codelingo.yaml file
    tenets:
      - name: example-tenet
        flows:
          codelingo/{{.Name}}:
        query:
          import codelingo/ast/go

          @{{.Name}} {{.DecoratorExample}}
          go.func_decl(depth := any)


`[1:]

var INFOTMP = `
  FLOW APP INFO
  =============
   
   Version
   -------

    {{.Version}} - last compiled {{.Compiled}}
   {{if len .Authors}}
   AUTHOR
   ------
   {{range .Authors}}{{ . }}{{end}}
   {{end}}

`[1:]
