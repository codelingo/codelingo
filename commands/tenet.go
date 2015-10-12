package commands

import (
	"fmt"

	"strings"

	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/commands/review"
	"github.com/lingo-reviews/lingo/drivers"
	"github.com/lingo-reviews/lingo/tenet"
)

// TenetCMD is a fallthrough CMD which treats command as the tenet name and
// passes through any arguments to the tenet.
func TenetCMD(ctx *cli.Context, command string) {
	// TODO(matt) read about bash completion on
	// https://github.com/codegangsta/cli. Is there a nice way that we could
	// add bash completion for tenet names (as they'll be long and clumsy).
	t, err := tenet.New(command)
	if err != nil {
		oserrf("command or tenet not found: %s", err.Error())
		return
	}

	// Create and initialise a driver
	driver, err := drivers.New(t.Driver, ctx)
	if err != nil {
		oserrf(err.Error())
		return
	}

	method := "Help"
	args := ctx.Args()
	if len(args[1:]) > 0 {
		method = methodFromArg(args[1])
	}

	switch method {
	case "Review":
		reviewResult, err := driver.Review(t, args[2:]...)
		if err != nil {
			oserrf("error running method %q, %s", method, err.Error())
			return
		}
		for _, i := range reviewResult.Issues {
			fmt.Println(review.FormatPlainText(i))
		}
		for _, e := range reviewResult.Errs {
			fmt.Println(e)
		}
	case "Help":
		text, err := driver.Help(t, args[2:]...)
		if err != nil {
			oserrf("error running method %q, %s", method, err.Error())
			return
		}
		fmt.Println(text)
	case "Version":
		text, err := driver.Version(t)
		if err != nil {
			oserrf("error running method %q, %s", method, err.Error())
			return
		}
		fmt.Println(text)
	case "Debug":
		fmt.Println(driver.Debug(t, args...))
	case "CommentSet":
		commSet, err := driver.CommentSet(t)
		if err != nil {
			oserrf("error running method %q, %s", method, err.Error())
			return
		}
		fmt.Println("\nThis command is provided for debugging purposes only. The following is a dump of the comments in the tenet's CommentSet.\n")
		for _, c := range commSet.Comments {
			fmt.Printf("\n%#v\n", c)
		}
	default:
		oserrf("tenet does not have method %q", method)
		return
	}
}

func methodFromArg(arg string) (method string) {
	method = strings.Title(arg)
	return strings.Replace(method, "-", "", -1)
}
