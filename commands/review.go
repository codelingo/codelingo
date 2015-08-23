package commands

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"

	t "github.com/lingo-reviews/dev/tenet"
	"github.com/lingo-reviews/lingo/commands/review"
	"github.com/lingo-reviews/lingo/tenet"
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
		cli.Float64Flag{
			Name:  "min-confidence",
			Value: 0,
			Usage: "the minimum confidence an issue needs to be included",
		},
		cli.IntFlag{
			Name:  "wait",
			Value: 20,
			Usage: "how long to wait, in seconds, for a tenet to finish.",
		},
		cli.StringFlag{
			Name:   "output",
			Value:  "cli",
			Usage:  "file path to write the output to. By default, output will be printed to the CLI",
			EnvVar: "LINGO-OUTPUT",
		},
		cli.StringFlag{
			Name:   "output-fmt",
			Value:  "plain-text",
			Usage:  "json, json-pretty, yaml, toml or plain-text. If an output-template is set, it takes precedence",
			EnvVar: "LINGO-OUTPUT-FMT",
		},
		cli.StringFlag{
			// TODO(waigani) implement. We could make output-fmt fall-through to check for custom template?
			Name:   "output-template",
			Value:  "",
			Usage:  "a template for the output format",
			EnvVar: "LINGO-OUTPUT-TEMPLATE",
		},
		cli.BoolFlag{
			Name:   "diff",
			Usage:  "only report issues found in unstaged, uncommited work",
			EnvVar: "LINGO-DIFF",
		},
		cli.BoolFlag{
			Name:   "no-user-confirm",
			Usage:  "turns off the default behaviour of stepping through each issue found and asking the user to confirm that it is an issue.",
			EnvVar: "LINGO-NO-USER-CONFIRM",
		},
	},
	Action: reviewAction,
}

func reviewAction(c *cli.Context) {
	f := "/source/file.go"
	p := c.GlobalString("repo-path")
	basePath, err := filepath.Abs(p)
	if err != nil {
		panic(err)
	}
	f = strings.Replace(f, "/source", basePath, 1)

	ts := tenets(c)
	// setup a chan of results.
	results := make(chan *tenet.ReviewResult, len(ts))
	var wg sync.WaitGroup
	wg.Add(len(ts))
	// wait for all results to come in before closing the chan.
	go func() {
		wg.Wait()
		close(results)
	}()

	for _, tn := range ts {
		go func() {
			defer wg.Done()

			err := tn.DockerInit()
			if err != nil {
				oserrf(err.Error())
				return
			}

			// TODO(waigani)
			// - no args should recursively review all files in pwd.
			// - --diff should drop any file not in the diff.
			args := c.Args()
			if opts := c.String("options"); opts != "" {
				args = append([]string{"--options", opts}, args...)
			}
			reviewResult, err := tn.Review(args...)
			if err != nil {
				oserrf("error running review %s", err.Error())
				return
			}
			// TODO(waigani) we can be smarter here. Pipe individual issues
			// from tenet to chan. Use fan-in pattern:
			// https://blog.golang.org/pipelines
			results <- reviewResult
		}()
	}

	r := allResults(c, results)

	if len(r.errors) > 0 {
		fmt.Println("The following errors were encounted:")
		for _, err := range r.errors {
			fmt.Printf("%v\n", err)
		}

		fmt.Println("Do you still wish to output the found issues? [y]es [N]o")

		var options string
		fmt.Print("\n[o]pen [d]iscard [K]eep:")
		fmt.Scanln(&options)

		switch options {
		case "y", "Y", "yes":
		default:
			return
		}

	}

	outputFmt := review.OutputFormat(c.String("output-fmt"))
	output := review.Output(outputFmt, c.String("output"), r.issues)
	fmt.Println(output)
}

type result struct {
	issues []*t.Issue
	errors []error
}

// TODO(waigani) TECHDEBT if diff is true, we only report the issues found
// within the diff, even though results contains all issues in the target
// file(s). Yes, this is just stupid. We need to pass the file diff boundaries
// to the tenets, it is then the tenet's responsibility to only analyse those
// nodes/lines within the diff.

// allResults returns all the issues all the tenets found.
func allResults(c *cli.Context, results chan *tenet.ReviewResult) result {
	issues := make(chan *t.Issue)
	tenetErrs := make(chan string)

	var wg sync.WaitGroup

	wait := time.Duration(int64(c.Int("wait"))) * time.Second
	var errs []error
l:
	for {
		select {
		case r, ok := <-results:
			if !ok {
				break l
			}
			wg.Add(len(r.Issues))
			wg.Add(len(r.Errs))
			go func() {
				for _, i := range r.Issues {
					issues <- i
					wg.Done()
				}
			}()

			go func() {
				for _, e := range r.Errs {
					tenetErrs <- e
					wg.Done()
				}
			}()
		case <-time.After(wait):
			msg := "timed out, the following tenet(s) did not run:"
			select {
			case r := <-results:
				msg += " " + r.TenetName
			default:
				errs = append(errs, errors.New(msg))
			}
		}
	}

	go func() {
		wg.Wait()
		close(issues)
		close(tenetErrs)
	}()

	cfm := review.NewConfirmer(c)

	var confirmedIssues []*t.Issue
	issuesClosed, errsClosed := false, false

	for {
		if issuesClosed && errsClosed {
			break
		}
		select {
		case issue, ok := <-issues:
			if !ok {
				issuesClosed = true
				continue
			}
			if cfm.Confirm(0, issue) {
				confirmedIssues = append(confirmedIssues, issue)
			}
		case errMsg, ok := <-tenetErrs:
			if !ok {
				errsClosed = true
				continue
			}
			errs = append(errs, errors.New(errMsg))
		case <-time.After(wait):
			msg := "timed out"
			errs = append(errs, errors.Errorf(msg))
		}
	}

	return result{confirmedIssues, errs}
}
