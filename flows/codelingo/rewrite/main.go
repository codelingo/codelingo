package main

import (
	"github.com/codelingo/codelingo/flows/codelingo/rewrite/rewrite"

	flowutil "github.com/codelingo/codelingo/sdk/flow"
	"github.com/codelingo/lingo/app/util"
	"github.com/juju/errors"
)

func main() {
	fRunner := flowutil.NewFlow(rewrite.CLIApp, rewrite.DecoratorApp)
	resultc, errc := fRunner.Run()
	var results []*flowutil.DecoratedResult

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

			results = append(results, result)
		}
		if resultc == nil && errc == nil {
			break l
		}
	}
	if hasErred {
		return
	}

	if err := rewrite.Write(results); err != nil {
		util.Logger.Debugw("Rewrite Flow", "err_stack", errors.ErrorStack(err))
		util.FatalOSErr(err)
		return
	}
}
