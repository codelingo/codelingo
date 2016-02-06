package commands

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/lingo-reviews/lingo/commands/common"
	"github.com/lingo-reviews/lingo/commands/review"
	"github.com/lingo-reviews/lingo/commands/review/service"
)

var ReviewCMD = cli.Command{
	Name:  "review",
	Usage: "review code following tenets in .lingo",
	Description: `

Review all files found in pwd, following tenets in .lingo of pwd or parent directory:
	"lingo review"

Review all files found in pwd, with two speific tenets:
	"lingo review \
	lingoreviews/space-after-forward-slash \
	lingoreviews/unused-args"

	This command ignores any tenets in any .lingo files.

`[1:],
	Subcommands: service.Services,
	Flags:       review.Flags,
	Action:      reviewAction,
}

func reviewAction(ctx *cli.Context) {
	opts := review.Options{
		Files:      ctx.Args(),
		Diff:       ctx.Bool("diff"),
		SaveToFile: ctx.String("save"),
		KeepAll:    ctx.Bool("keep-all"),
	}
	issues, err := review.Review(opts)
	if err != nil {
		common.OSErrf(err.Error())
		return
	}

	// TODO(waigani) I'm not happy that we're doing this here, can't find a
	// better place though.
	saveToFile := ctx.String("save")
	if saveToFile != "" {
		err := review.Save(saveToFile, issues)
		if err != nil {
			common.OSErrf("could not save to file: %s", err.Error())
			return
		}
		fmt.Printf("review saved to %s\n", saveToFile)
	}

	// TODO(waigani) make more informative
	// TODO(waigani) if !ctx.String("quiet")
	fmt.Printf("Done! Found %d issues \n", len(issues))
}
