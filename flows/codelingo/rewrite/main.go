package main

import (
	"github.com/codelingo/codelingo/flows/codelingo/rewrite/rewrite"

	flowutil "github.com/codelingo/codelingo/sdk/flow"
	"github.com/codelingo/lingo/app/util"
	"github.com/juju/errors"
)

func main() {
	fRunner := flowutil.NewFlow(rewrite.CLIApp, rewrite.DecoratorApp)
	resultc, err := fRunner.Run()
	if err != nil {
		util.Logger.Debugw("Rewrite Flow", "err_stack", errors.ErrorStack(err))
		util.FatalOSErr(err)
		return
	}

	var results []*flowutil.DecoratedResult
	for result := range resultc {
		results = append(results, result)
	}

	if err := rewrite.Write(results); err != nil {
		util.Logger.Debugw("Rewrite Flow", "err_stack", errors.ErrorStack(err))
		util.FatalOSErr(err)
		return
	}
}
