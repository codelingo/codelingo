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

	fRunner := flowutil.NewFlow(review.CLIApp, review.DecoratorApp)
	resultc, errc := fRunner.Run()

	// TODO(waigani) capture false positives in the confirmer.
	var results []*review.ReportStrt

	var hasErred bool

l:
	for {
		select {
		case err, ok := <-errc:
			if !ok {
				errc = nil
				break
			}

			util.Logger.Debugw("Rewrite Flow", "err_stack", errors.ErrorStack(err))
			util.FatalOSErr(err)
			hasErred = true
		case result, ok := <-resultc:
			if !ok {
				resultc = nil
				break
			}

			issue := result.Payload.(*flow.Issue)
			results = append(results, &review.ReportStrt{
				Name:     issue.Name,
				Comment:  issue.Comment,
				Filename: issue.Position.Start.Filename,
				Line:     int(issue.Position.Start.Line),
				Snippet:  issue.CtxBefore + "\n" + issue.LineText + "\n" + issue.CtxAfter,
			})

		}
		if resultc == nil && errc == nil {
			break l
		}
	}

	if hasErred {
		return
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

	if cliCtx.IsSet("output") {
		msg, err := review.MakeReport(cliCtx, results)
		if err != nil {
			util.Logger.Debugw("Review Flow", "err_stack", errors.ErrorStack(err))
			util.FatalOSErr(err)
			return
		}
		fmt.Println(msg)
	}

}
