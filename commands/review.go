package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/lingo-reviews/dev/api"
	"google.golang.org/grpc/grpclog"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/waigani/diffparser"

	"github.com/lingo-reviews/lingo/commands/review"
	"github.com/lingo-reviews/lingo/tenet"
	// TODO: Avoid driver import
	"github.com/lingo-reviews/lingo/tenet/driver"
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
			Value:  "none",
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
			Name:   "keep-all",
			Usage:  "turns off the default behaviour of stepping through each issue found and asking the user to confirm that it is an issue.",
			EnvVar: "LINGO-KEEP-ALL",
		},
	},
	Action: reviewAction,
}

// TODO(waigani) this is ~300 line long func! Break it up.
func reviewAction(ctx *cli.Context) {
	var diff *diffparser.Diff
	var err error
	rQueue := make(map[*config][]TenetConfig)
	totalTenets := 0
	files := ctx.Args()

	// create new diff to filter issues by.
	if ctx.Bool("diff") {
		diff, err = diffparser.Parse(rawDiff())
		if err != nil {
			oserrf(err.Error())
			return
		}
	}

	// if no files are named and we are diffig, add all files in diff.
	if len(files) == 0 && diff != nil {
		for _, f := range diff.Files {
			// TODO(waigani) DEMOWARE make "tenet.toml" a cfg var. We should
			// support reviewing the cfg also, right now it errors out.
			if f.Mode != diffparser.DELETED && !strings.Contains(f.NewName, "tenet.toml") {
				files = append(files, f.NewName)
			}
		}
	}

	// Get this first as it might fail, we want to avoid all other work in that case.
	cfm, err := review.NewConfirmer(ctx, diff)
	if err != nil {
		oserrf(err.Error())
		return
	}

	if len(files) > 0 {
		for _, file := range files {
			cfgPath := path.Join(path.Dir(file), defaultTenetCfgPath)
			cfg, err := buildConfig(cfgPath, CascadeUp)
			if err != nil {
				oserrf(err.Error())
				return
			}

			for _, tenetCfg := range cfg.AllTenets() {
				totalTenets++
				rQueue[cfg] = append(rQueue[cfg], tenetCfg)
			}
		}
	} else {
		// TODO: Check for dirs amongst files
		rQueue, totalTenets, err = reviewQueue(".")
		if err != nil {
			oserrf(err.Error())
			return
		}
	}

	// TODO(waigani) I'm not sure this is the right place to be merging in CLI options. tenet.New takes ctx, so why not there?
	commandOptions, err := parseOptions(ctx)
	if err != nil {
		oserrf(err.Error())
		return
	}

	// setup a chan of results.
	resultsc := make(chan *result, totalTenets)
	var wg sync.WaitGroup
	wg.Add(totalTenets)
	grpclog.Println("totalTenets", totalTenets)
	// wait for all results to come in before closing the chan.
	go func() {
		wg.Wait()
		close(resultsc)
	}()

	// A map of all matching files for dir and extension.
	fileMatches := map[string][]string{}

	for cfg, tenetCfgs := range rQueue {
		// TODO(waigani) wait for all tenets to review all files.
		for _, tenetCfg := range tenetCfgs {

			// setup results for this tenet.
			r := &result{
				tenetName: tenetCfg.Name,

				// Setting a buffer will allow the tenet to continue
				// to find and load issues while the user confirms.
				issuesc: make(chan *api.Issue, 5),
			}

			// Wrap in a closure so we can process any error with the one defer.
			func(tenetCfg TenetConfig, r *result) (err error) {
				defer func() {
					// Finish early and close the issues chan if we've returned on an error.
					if err != nil {
						r.err = errors.Annotatef(err, "tenet %q errored while reviewing", r.tenetName)
						grpclog.Print(r.err)
						close(r.issuesc)

						resultsc <- r
					}
				}()

				// Merge in options from command line.
				for k, v := range commandOptions[tenetCfg.Name] {
					tenetCfg.Options[k] = v
				}

				var tn tenet.Tenet
				tn, err = tenet.New(ctx, &driver.Base{
					Name:          tenetCfg.Name,
					Driver:        tenetCfg.Driver,
					Registry:      tenetCfg.Registry,
					Tag:           tenetCfg.Tag,
					ConfigOptions: tenetCfg.Options,
				})
				if err != nil {
					return
				}

				// Start the review.
				go func(tn tenet.Tenet) {
					// get the tenet service
					var s tenet.TenetService
					s, err = tn.Service()
					if err != nil {
						return
					}
					if err = s.Start(); err != nil {
						return
					}

					// Build up the file names for this tenet to review.
					filesc := make(chan string)

					// Start listening for files to review.
					go func(s tenet.TenetService, filesc chan string) {
						// // TODO(waigani) do we care about stop errs?
						defer s.Stop()
						if err = s.Review(filesc, r.issuesc); err != nil {
							return
						}
						wg.Done()
					}(s, filesc)

					// Send files to review.
					go func(s tenet.TenetService, filesc chan string) {
						var lang string
						lang, err = s.Language()
						if err != nil {
							return
						}

						// Find files for this tenet.
						regex, glob := fileExtFilterForLang(lang)
						globSearch := path.Join(cfg.buildRoot, glob)

						if fileMatches[globSearch] == nil {
							// Use named files if passed in.
							var fNames []string
							if len(files) > 0 {
								// Add only those files that this tenet is interested in.
								for _, fName := range files {
									if m, err := regexp.MatchString(regex, fName); !m {
										if err != nil {
											// TODO(waigani) log msg here
										}
										continue
									}
									fNames = append(fNames, fName)
								}
							} else {
								fNames, err = filepath.Glob(globSearch)
								if err != nil { // Non-fatal.
									// TODO(waigani) log: no files for this tenet to review.
								}
							}

							// Do some final checks on each file.
							// If this is a diff check and file is not in diff, don't review it.
							// Ensure the file can be opened without error and is
							// not a directory.
							for _, f := range fNames {

								// !!!! TECHDEBT DEMOWARE MATT REMIND ME
								// WE NEED TO FIX THIS BEFORE ALPHA !!!!. TODO(waigani) don't add files not in
								// diff.  The problem is, if the context gets filled on files not in the diff,
								// the tenet could stop before getting to the diffed files.
								// if c.Bool("diff"){
								// 	if file is not in diff, continue
								// }

								file, err := os.Open(f)
								if err != nil { // Non-fatal
									file.Close()
									continue
								}
								if fi, err := file.Stat(); err == nil && !fi.IsDir() {
									// fmt.Println("adding", f) // TODO: put behind a debug flag
									fileMatches[globSearch] = append(fileMatches[globSearch], f)
								} else {
									// TODO(waigani) log here.
								}
								file.Close()
							}

							// sort the order of files so the context set by the tenet will be correct.
							sort.Strings(fileMatches[globSearch])
						}

						// Push files to be reviewed by this tenet onto the chan
						// and close the chan once done.
						for _, fName := range fileMatches[globSearch] {
							filesc <- fName
						}
						close(filesc)
					}(s, filesc)

					// If there was an error, the defer func above will register
					// that error and finish up for us.
					if err == nil {
						// fan in our review result. result contains an issue chan
						// which will be listened on until closed.
						resultsc <- r
					}
				}(tn)
				return
			}(tenetCfg, r)
		}
	}
	grpclog.Println("colating results")
	issues, errs := allResults(ctx, cfm, resultsc)
	grpclog.Println("review done")
	// Print the final output to the user.
	if len(errs) > 0 {
		fmt.Println("The following errors were encounted:")
		for _, err := range errs {
			fmt.Printf("%v\n", err)
		}

		var options string
		fmt.Println("Do you still wish to output the found issues? [y]es [N]o")
		fmt.Scanln(&options)

		switch options {
		case "y", "Y", "yes":
		default:
			return
		}
	}

	// Even if there are no issues, we still might need to show output.
	outputFmt := review.OutputFormat(ctx.String("output-fmt"))
	if outputFmt != "none" {
		output := review.Output(outputFmt, ctx.String("output"), issues)
		fmt.Print(output)
	}
}

type result struct {
	tenetName string
	issuesc   chan *api.Issue
	err       error
}

// TODO(waigani) TECHDEBT if diff is true, we only report the issues found
// within the diff, even though results contains all issues in the target
// file(s). We need to pass the diff  hunk boundaries to the tenets, it is then
// the tenet's responsibility to only analyse those nodes/lines within the
// hunks.

// allResults returns all the issues all the tenets found.
func allResults(c *cli.Context, cfm *review.IssueConfirmer, resultsc chan *result) ([]*api.Issue, []error) {
	allIssues := make(chan *api.Issue)
	var errs []error
	// TODO(waigani) chan of comment context.

	// if --keep-all, we sort after all issues are found and then apply the
	// comment context. Otherwise, we use the  as the issues come in.

	go func(allIssues chan *api.Issue) {
		for r := range resultsc {
			if r.err != nil {
				errs = append(errs, errors.Annotatef(r.err, "tenet %q", r.tenetName))
			}

			go func(r *result) {
				for i := range r.issuesc {
					allIssues <- i
				}
			}(r)
		}
		close(allIssues)
	}(allIssues)

	var confirmedIssues []*api.Issue
	for {
		select {
		case issue, ok := <-allIssues:
			if !ok {
				break
			}

			if cfm.Confirm(0, issue) {
				confirmedIssues = append(confirmedIssues, issue)
			}
		case <-time.After(20 * time.Second):
			errs = append(errs, errors.New("timed out"))
		}
	}

	return confirmedIssues, errs
}

// TODO(waigani) this just reads unstaged changes from git in pwd. Change diff
// from a flag to a sub command which pipes files to git diff.
func rawDiff() string {
	c := exec.Command("git", "reset")
	c.Run()
	c = exec.Command("git", "add", "-N", ".") // this includes new files in diff
	c.Run()

	var stdout bytes.Buffer
	c = exec.Command("git", "diff")
	c.Stdout = &stdout
	// c.Stderr = &stderr
	c.Run()
	diff := string(stdout.Bytes())

	c = exec.Command("git", "reset")
	c.Run()

	return diff
}
