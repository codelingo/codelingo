package commands

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"

	"github.com/codegangsta/cli"
	"github.com/waigani/diffparser"
	tomb "gopkg.in/tomb.v1"

	"github.com/lingo-reviews/dev/api"
	"github.com/lingo-reviews/dev/tenet/log"

	"github.com/lingo-reviews/lingo/commands/review"
)

var ReviewCMD = cli.Command{
	Name:  "review",
	Usage: "review code following tenets in .lingo",
	Description: `

Review all files found in pwd, following tenets in .lingo of pwd or parent directory:
	"lingo review"

Review all files found in pwd, with two speific tenets:
	"lingo review \
	lingoreviews/space-after-forward-slash \
	lingoreviews/unused-args"

	This command ignores any tenets in any .lingo files.

`[1:],
	Flags: []cli.Flag{
		cli.StringFlag{
			// TODO(waigani) interactively set options for tenet.
			Name:  "options",
			Usage: "serialized JSON options from .lingo",
		},
		cli.StringFlag{
			Name:   "output",
			Value:  "cli",
			Usage:  "file path to write the output to. By default, output will be printed to the CLI",
			EnvVar: "LINGO-OUTPUT",
		},
		cli.StringFlag{
			Name:  "output-fmt",
			Value: "none",
			// TODO(waigani) support yaml toml. Also: if an output-template is set, it takes precedence.
			Usage:  "json or json-pretty",
			EnvVar: "LINGO-OUTPUT-FMT",
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
		cli.BoolFlag{
			Name:   "find-all",
			Usage:  "raise every issue tenets find",
			EnvVar: "LINGO-KEEP-ALL",
		},
		// TODO(waigani) implement
		// cli.IntFlag{
		// 	Name:  "wait",
		// 	Value: 20,
		// 	Usage: "how long to wait, in seconds, for a tenet to finish.",
		// },
		// cli.StringFlag{
		// 	// TODO(waigani) implement. We could make output-fmt fall-through to check for custom template?
		// 	Name:   "output-template",
		// 	Value:  "",
		// 	Usage:  "a template for the output format",
		// 	EnvVar: "LINGO-OUTPUT-TEMPLATE",
		// },
		// TODO(waigani) implement
		// cli.BoolFlag{
		// 	Name:  "watch",
		// 	Usage: "watch for changes and review any changed files",
		// },
	},
	Action: reviewAction,
}

type cfgMap struct {
	path  string
	cfg   *config
	dirs  []string
	files []string
}

// readCfgs converts work done up front (finding configs) into a concurrent pipeline
// stage (reading configs).
func readCfgs(cfgList []cfgMap, errc chan error) <-chan cfgMap {
	var err error
	out := make(chan cfgMap)
	go func() {
		for _, m := range cfgList {
			m.cfg, err = buildConfig(m.path, CascadeUp)
			if err != nil {
				errc <- err
				continue
			}
			out <- m
		}
		close(out)
	}()
	return out
}

type tenetReview struct {
	configHash string
	filesc     chan string
	issuesc    chan *api.Issue
	info       *api.Info
	issuesWG   *sync.WaitGroup
	filesTM    *tomb.Tomb
}

// TODO(waigani) use configHash for docker containers. Only remove the container at the end of review.
// TODO(waigani) add a buffer. We want lingo to be ahead of the user, but not to review the world in the background.
// If the user has one issue up on their screen, very little cpu / mem should be being used.

var bufferFullERR = errors.New("buffer full")

// returns a chan of tenet reviews and a cancel chan that blocks until the user cancels.
func reviewQueue(ctx *cli.Context, mappings <-chan cfgMap, errc chan error) (<-chan *tenetReview, chan struct{}) {
	reviews := make(map[string]*tenetReview)
	reviewChannel := make(chan *tenetReview)
	cleanupWG := &sync.WaitGroup{}

	// setup a cancel exit path.
	cancelc := make(chan os.Signal, 1)
	signal.Notify(cancelc, os.Interrupt)
	signal.Notify(cancelc, syscall.SIGTERM)
	cancelledc := make(chan struct{})
	cancelled := func() bool {
		select {
		case _, ok := <-cancelledc:
			if ok {
				close(cancelc)
			}
			return true
		default:
			return false
		}
	}

	// Kill all open tenets on cancel.
	go func() {
		var i int
		for {
			<-cancelc
			if i > 0 {
				// on the second exit, just do it.
				fmt.Print("failed.\nSome docker containers may still be running.")
				os.Exit(1)
			}
			i++
			go func() {
				// TODO(waigani) add progress bar here
				fmt.Print("\ncleaning up tenets ... ")

				// Anything waiting on the cancelled chan will now fire.
				close(cancelledc)

				// Wait for all tenets to be cleaned up.
				cleanupWG.Wait()

				// say bye.
				fmt.Println("done.")
				os.Exit(1)
			}()
		}
	}()

	// TODO(waigani) reenable buffering to:
	// 1. Allow found tenets to keep running.
	// 2. Stop building new tenets until there is room in the buffer.
	// TODO(waigani) make cfg vars
	// buffLimit := 3
	// if ctx.Bool("keep-all") {
	// 	buffLimit = 100
	// }
	// buff := util.NewBuffer(buffLimit, cancelledc)

	go func() {
		for m := range mappings {
			// Glob all the files in the associated directories for this config, and assign to each tenet by hash
			for _, tc := range m.cfg.AllTenets() {

				if cancelled() {
					// empty dirs to stop feeding tenet reviews in progress.
					m.dirs = []string{}
					return
				}

				// Open the tenet service if we haven't seen this config before.
				configHash := tc.hash()
				r, found := reviews[configHash]
				if !found {

					// Don't build a new tenet until there is room in the buffer.
					// Found tenets will keep running until they are not fed files for 5 seconds.
					// WaitRoom will not block if we get a cancel signal.
					// buff.WaitRoom()

					tn, err := newTenet(ctx, tc)
					if err != nil {
						errc <- err
						continue
					}
					// Note: service should not be called outside this if block.
					service, err := tn.OpenService()
					if err != nil {
						errc <- err
						continue
					}

					info, err := service.Info()
					if err != nil {
						errc <- err
						continue
					}

					r = &tenetReview{
						configHash: configHash,
						filesc:     make(chan string),
						issuesc:    make(chan *api.Issue),
						info:       info,
						issuesWG:   &sync.WaitGroup{},
						filesTM:    &tomb.Tomb{},
					}
					reviews[configHash] = r

					// Setup the takedown of this review.
					r.issuesWG.Add(1)
					cleanupWG.Add(1)
					// buff.Add(1)

					go func(r *tenetReview) {
						// The following fires when:
						select {
						// 1. all files have been sent or timed out
						// 2. the tenet buffer is full
						case <-r.filesTM.Dying():
							// 3. lingo has been stopped
						case <-cancelledc:
						}

						// make room for another tenet to start and ensure
						// that any configHash's matching this one will have
						// to start a new tenet instance.
						delete(reviews, configHash)
						// buff.Add(-1)

						// signal to the tenet that no more files are coming.

						close(r.filesc)

						// wait for the tenet to signal to us that it's finished it's review.
						r.issuesWG.Wait()

						// we can now safely close the backing service.
						if err := service.Close(); err != nil {
							log.Println("ERROR closing sevice:", err)
						}

						log.Println("cleanup done")
						cleanupWG.Done()
					}(r)

					// Make sure we're ready to handle results before we start
					// the review.
					reviewChannel <- r

					// Start this tenet's review.
					service.Review(r.filesc, r.issuesc, r.filesTM)
				}

				regexPattern, globPattern := fileExtFilterForLang(r.info.Language)
				for _, d := range m.dirs {
					files, err := filepath.Glob(path.Join(d, globPattern))
					if err != nil {
						// Non-fatal
						log.Printf("Error reading files in %s: %v\n", d, err)
					}

				l:
					for i, f := range files {
						select {
						case <-cancelledc:
							log.Println("user cancelled, dropping files.")
						case <-r.filesTM.Dying():
							dropped := len(files) - i
							log.Print("WARNING a tenet review timed out waiting for files to be sent. %d files dropped", dropped)
							break l
						case r.filesc <- f:
						}
					}
				z:
					for i, f := range m.files {
						if matches, err := regexp.MatchString(regexPattern, f); !matches {
							if err != nil {
								log.Println("error in regex: ", regexPattern)
							}
							continue
						}

						select {
						case <-cancelledc:
							log.Println("user cancelled, dropping files.")
						case <-r.filesTM.Dying():
							dropped := len(m.files) - i
							log.Print("WARNING a tenet review timed out waiting for files to be sent. %d files dropped", dropped)
							break z
						case r.filesc <- f:
						}
					}

				}

			}
		}

		for _, r := range reviews {

			// this says all files have been sent. For this review.
			r.filesTM.Done()
		}

		// wait for all tenets to be cleaned up.
		cleanupWG.Wait()

		// Closing this chan will start the wind down to end the lingo
		// process.
		close(reviewChannel)
	}()

	return reviewChannel, cancelledc
}

func reviewAction(ctx *cli.Context) {

	// TODO: file args input, as files and dirs
	var err error
	var diff *diffparser.Diff

	// create new diff to filter issues by
	if ctx.Bool("diff") {
		diff, err = diffparser.Parse(rawDiff())
		if err != nil {
			oserrf(err.Error())
			return
		}
	}

	// Get this first as it might fail, we want to avoid all other work in that case
	cfm, err := review.NewConfirmer(ctx, diff)
	if err != nil {
		oserrf(err.Error())
		return
	}

	fileArgs := ctx.Args()
	if diff != nil {
		// if we are diffing, add all files in diff
		for _, f := range diff.Files {
			if f.Mode != diffparser.DELETED {
				fileArgs = append(fileArgs, f.NewName)
			}
		}
	} else if len(fileArgs) == 0 {
		fileArgs = []string{"."}
	}

	// Receiver for errors that can occur during pipeline stages
	errc := make(chan error)
	errors := []error{}
	// Just collect errors during review - show them to the user at the end
	go func() {
		for err := range errc {
			errors = append(errors, err)
		}
	}()

	// Map of project config filenames -> directories they control
	cfgList := []cfgMap{}
	// TODO: This loop is now pipelinable too, if we need to further reduce time-to-first-review
	for _, f := range fileArgs {
		// Specifically asking for a file that can't be found/read is fatal
		file, err := os.Open(f)
		if err != nil {
			oserrf(err.Error())
			return
		}
		fi, err := file.Stat()
		if err != nil {
			oserrf(err.Error())
			return
		}

		if fi.IsDir() {
			filepath.Walk(f, func(relPath string, info os.FileInfo, err error) error {
				if info.IsDir() {
					// Ignore folders beginning with '.', except search root
					// TODO: Flag to turn this behaviour off
					if len(relPath) > 1 && info.Name()[0] == '.' {
						return filepath.SkipDir
					}
					// TODO: Faster technique for finding cfgPath taking advantage of Walk's depth-first search
					//       This implementation recurses upwards for each found dir
					cfgPath, _ := tenetCfgPathRecusive(path.Join(relPath, defaultTenetCfgPath))
					cfgList = append(cfgList, cfgMap{
						path:  cfgPath,
						cfg:   nil,
						dirs:  []string{relPath},
						files: []string{},
					})
				}
				return nil
			})
		} else {
			cfgPath, _ := tenetCfgPathRecusive(path.Join(filepath.Dir(f), defaultTenetCfgPath))
			cfgList = append(cfgList, cfgMap{
				path:  cfgPath,
				cfg:   nil,
				dirs:  []string{},
				files: []string{f},
			})
		}
	}

	// Use a channel to read configs with directory mapping
	configDirs := readCfgs(cfgList, errc)

	rc, cancelledc := reviewQueue(ctx, configDirs, errc)
	var count int

	keptIssuesc := make(chan *api.Issue)

	// collectedIssues has a huge buffer so we can store all the found issues,
	// allowing the tenet instances to be stopped. If this buffer is filled,
	// tenets will not be stopped. They will hang around until there is room
	// to offload their issues.
	collectedIssuesc := make(chan *api.Issue, 1000000)
	allIssuesWG := &sync.WaitGroup{}

	// Confirm found issues
	waitToConfirmc := make(chan struct{})
	if ctx.Bool("keep-all") {
		close(waitToConfirmc)
	}

	go func() {
		for i := range collectedIssuesc {
			if cfm.Confirm(0, i) {
				// Don't block on send. In the case of --keep-all
				// with no output, we just show count and have no
				// need for issues.
				select {
				case <-waitToConfirmc:
				case keptIssuesc <- i:
				}
			}
		}
		close(keptIssuesc)
	}()

z:
	for {
		select {
		case r, open := <-rc:
			if !open && r == nil {
				break z
			}
			allIssuesWG.Add(1)
			go func(r *tenetReview) {
				defer r.issuesWG.Done()

			l:
				for {
					select {
					case i, ok := <-r.issuesc:
						if !ok && i == nil {
							allIssuesWG.Done()
							break l
						}
						count++
						select {
						case <-cancelledc:
						case collectedIssuesc <- i:
						}

					}
				}

			}(r)
		}
	}

	// Wait for all issues to be read.
	allIssuesWG.Wait()

	// then close our collection chan.
	close(collectedIssuesc)

	var issues []*api.Issue
	for i := range keptIssuesc {
		issues = append(issues, i)
	}

	outputFmt := review.OutputFormat(ctx.String("output-fmt"))
	if outputFmt != "none" {
		output := review.Output(outputFmt, ctx.String("output"), issues)
		fmt.Print(output)
	} else {

		// TODO(waigani) make more informative
		// TODO(waigani) if !ctx.String("quiet")
		fmt.Printf("Done! Found %d issues \n", count)
	}

	close(errc)

	// Print errors if any occured
	if len(errors) > 0 {
		fmt.Println("The following errors were encounted:")
		for _, err := range errors {
			fmt.Printf("%v\n", err)
		}

		if outputFmt != "none" {
			var options string
			fmt.Println("Do you still wish to output the found issues? [y]es [N]o")
			fmt.Scanln(&options)

			switch options {
			case "y", "Y", "yes":
			default:
				return
			}
		}
	}
}

// TODO(waigani) this just reads unstaged changes from git in pwd. Change diff
// from a flag to a sub command which pipes files to git diff.
func rawDiff() string {
	c := exec.Command("git", "reset")
	c.Run()
	c = exec.Command("git", "add", "-N", ".") // this includes new files in diff
	c.Run()

	var stdout bytes.Buffer
	// TODO: Whilst --relative does get the correct files reviewed, corfirm/diffparser.Changed()
	//       is not receiving the correct paths and throws all diffs out
	c = exec.Command("git", "diff", "--relative")
	c.Stdout = &stdout
	// c.Stderr = &stderr
	c.Run()
	diff := string(stdout.Bytes())

	c = exec.Command("git", "reset")
	c.Run()

	return diff
}
