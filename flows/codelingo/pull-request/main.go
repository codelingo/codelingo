package flows

import (
	"context"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/codelingo/codelingo/flows/codelingo/review/review"
	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/rpc/flow"
	flowutil "github.com/codelingo/sdk/flow"
	"github.com/juju/errors"
)

var pullRequestCmd = cli.Command{
	Name:      "pull-request",
	ShortName: "pr",
	Usage:     "review a remote pull-request",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  util.LingoFile.String(),
			Usage: "A list of codelingo.yaml files to perform the review with. If the flag is not set, codelingo.yaml files are read from the branch being reviewed.",
		},
		// TODO(waigani) as this is a review sub-command, it should be able to use the
		// lingo-file flag from review.
		// cli.BoolFlag{
		// 	Name:  "all",
		// 	Usage: "review all files under all directories from pwd down",
		// },
	},
	Description: `
"$ lingo review pull-request https://github.com/codelingo/lingo/pull/1" will review all code in the diff between the pull request and it's base repository.
"$ lingo review pr https://github.com/codelingo/lingo/pull/1" will review all code in the diff between the pull request and it's base repository.
`[1:],
	// "$ lingo review" will review any unstaged changes from pwd down.
	// "$ lingo review [<filename>]" will review any unstaged changes in the named files.
	// "$ lingo review --all [<filename>]" will review all code in the named files.
	Action: reviewPullRequestAction,
}

func main() {
	if err := flowutil.Run(pullRequestCmd); err != nil {
		flowutil.HandleErr(err)
	}
}

func reviewPullRequestAction(ctx *cli.Context) {
	msg, err := reviewPullRequestCMD(ctx)
	if err != nil {
		// Debugging
		// print(errors.ErrorStack(err))
		util.FatalOSErr(err)
		return
	}
	fmt.Println(msg)
}

func reviewPullRequestCMD(cliCtx *cli.Context) (string, error) {
	if l := len(cliCtx.Args()); l != 1 {
		return "", errors.Errorf("expected one arg, got %d", l)
	}

	dotlingo, err := review.ReadDotLingo(cliCtx)
	if err != nil {
		return "", errors.Trace(err)
	}

	opts, err := review.ParsePR(cliCtx.Args()[0])
	if err != nil {
		return "", errors.Trace(err)
	}

	ctx, cancel := util.UserCancelContext(context.Background())
	issuec, errorc, err := review.RequestReview(ctx, &flow.ReviewRequest{
		Host:     opts.Host,
		Hostname: opts.HostName,
		// TODO (Junyu) separate it into two separate fields
		OwnerOrDepot: &flow.ReviewRequest_Owner{opts.Owner},
		Repo:         opts.Name,
		// Sha and patches are defined by the PR
		IsPullRequest: true,
		PullRequestID: int64(opts.PRID),
		Dotlingo:      dotlingo,
	})
	if err != nil {
		return "", errors.Trace(err)
	}

	issues, err := review.ConfirmIssues(cancel, issuec, errorc, cliCtx.Bool("keep-all"), cliCtx.String("save"))
	if err != nil {
		return "", errors.Trace(err)
	}

	// TODO: streaming back to the client, verify issues on the client side.
	if len(issues) == 0 {
		return "Done! No issues found.\n", nil
	}

	msg, err := review.MakeReport(issues, cliCtx.String("format"), cliCtx.String("save"))
	if err != nil {
		return "", errors.Trace(err)
	}

	fmt.Println(fmt.Printf("Done! Found %d issues \n", len(issues)))
	return msg, nil
}
