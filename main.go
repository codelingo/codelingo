package main

import (
	"fmt"
	"os"

	"github.com/juju/errors"
	"github.com/lingo-reviews/dev/tenet/log"
	"github.com/lingo-reviews/lingo/app"
)

func main() {
	err := app.New().Run(os.Args)
	if err != nil {
		if errors.Cause(err).Error() == "ui" {
			if e, ok := err.(*errors.Err); ok {
				log.Println(e.Underlying())
				fmt.Println(e.Underlying())
				os.Exit(1)
			}
		}
		panic(errors.ErrorStack(err))
	}
}
