package main

import (
	"os"

	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/app"
)

func main() {
	err := app.New().Run(os.Args)
	if err != nil {
		panic(errors.ErrorStack(err))
	}
}
