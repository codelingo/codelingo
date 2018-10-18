package main

import (
	"fmt"

	"github.com/codelingo/codelingo/flows/codelingo/review/review"
	flowutil "github.com/codelingo/codelingo/sdk/flow"
	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/rpc/flow"
	"github.com/juju/errors"
)

func main() {

	fRunner := flowutil.NewFlow(review.CliCMD, review.DecoratorCMD)
	resultc, err := fRunner.Run()
	if err != nil {
		util.Logger.Debugw("Review Flow", "err_stack", errors.ErrorStack(err))
		util.FatalOSErr(err)
		return
	}

	// TODO(waigani) capture false positives in the confirmer.
	var results []*review.ReportStrt
	for result := range resultc {

		issue := result.Payload.(*flow.Issue)
		results = append(results, &review.ReportStrt{
			Comment:  issue.Comment,
			Filename: issue.Position.Start.Filename,
			Line:     int(issue.Position.Start.Line),
			Snippet:  issue.CtxBefore + "\n" + issue.LineText + "\n" + issue.CtxAfter,
		})
	}

	if len(results) == 0 {
		fmt.Println("Done! No issues found.")
		return
	}

	fmt.Printf("Done! %d issues found.\n", len(results))

	cliCtx, err := fRunner.CliCtx()
	if err != nil {
		util.Logger.Debugw("Review Flow", "err_stack", errors.ErrorStack(err))
		util.FatalOSErr(err)
		return
	}

	msg, err := review.MakeReport(cliCtx, results)
	if err != nil {
		util.Logger.Debugw("Review Flow", "err_stack", errors.ErrorStack(err))
		util.FatalOSErr(err)
		return
	}

	fmt.Println(msg)

}
