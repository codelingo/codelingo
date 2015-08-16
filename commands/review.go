package commands

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var ReviewCMD = cli.Command{
	Name:  "review",
	Usage: "review code following tenets in tenet.toml",
	Description: `

Review all files found in pwd, following tenets in .lingo of pwd or parent directory:
	"lingo review"

Review all files found in pwd, with two speific tenets:
	"lingo review \
	lingoreviews/space-after-forward-slash \
	lingoreviews/unused-args"

	This command ignores any tenets in any tenet.toml files.

`[1:],
	Flags: []cli.Flag{
		cli.StringFlag{
			// TODO(waigani) interactively set options for tenet.
			Name:  "options",
			Usage: "serialized JSON options from tenet.toml",
		},
	},
	// TODO(waigani) add these flags
	// Flag{"output-type", "plain-text", "json, json-pretty, yaml, toml or plain-text. If an output-template is set, it takes precedence"},
	// Flag{"output-template", "", "a template for the output format"},
	// Flag{"diff", false, "only report issues found in unstaged, uncommited work"},
	// Flag{"interactive", false, "Step through each issue found and decide whether or not to remove it from the final output."},
	Action: reviewAction,
}

func reviewAction(c *cli.Context) {
	for _, t := range tenets(c) {
		err := t.DockerInit()
		if err != nil {
			oserrf(err.Error())
			return
		}

		args := c.Args()
		if opts := c.String("options"); opts != "" {
			args = append([]string{"--options", opts}, args...)
		}
		reviewResult, err := t.Review(args...)
		if err != nil {
			oserrf("error running review %s", err.Error())
			return
		}
		fmt.Println("tenet: ", t.String())
		for _, i := range reviewResult.Issues {
			// TODO(matt) currently formatting is in Comment func within
			// lingo-reviews/dev/tenet. Move the formatting to lingo.
			fmt.Println(i.String())
		}
		for _, e := range reviewResult.Errs {
			fmt.Println(e)
		}
	}
}
