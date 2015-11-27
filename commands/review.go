package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/waigani/diffparser"

	"github.com/lingo-reviews/dev/api"
	"github.com/lingo-reviews/dev/tenet/log"

	"github.com/lingo-reviews/lingo/commands/review"
	"github.com/lingo-reviews/lingo/tenet"
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

type pendingReview struct {
	configHash string
	filePath   string
}

type reviewStream struct {
	configHash string
	service    tenet.TenetService
	filesc     chan string
	issuesc    chan *api.Issue
}

// reviewQueue reads a channel of mappings and returns two output channels. The first receives new
// review objects that are generated when a new unique tenet config is encountered, the second is
// for each file/config pair to be reviewed and points back to a review object by hash.
func reviewQueue(ctx *cli.Context, mappings <-chan cfgMap, errc chan error) (<-chan reviewStream, <-chan pendingReview) {
	reviews := make(map[string]reviewStream)

	reviewChannel := make(chan reviewStream)
	pendingChannel := make(chan pendingReview)

	go func() {
		for m := range mappings {
			// Glob all the files in the associated directories for this config, and assign to each tenet by hash
			for _, tc := range m.cfg.AllTenets() {
				// Open the tenet service if we haven't seen this config before
				configHash := tc.hash()
				r, found := reviews[configHash]
				if !found {
					tn, err := newTenet(ctx, tc)
					if err != nil {
						errc <- err
						continue
					}
					service, err := tn.OpenService()
					if err != nil {
						errc <- err
						continue
					}
					r = reviewStream{
						configHash: configHash,
						service:    service,
						filesc:     make(chan string),
						issuesc:    make(chan *api.Issue, 5),
					}
					service.Review(r.filesc, r.issuesc)
					reviews[configHash] = r

					reviewChannel <- r
				}

				// TODO: Can cache this on reviewStream instead of asking for each file
				info, err := r.service.Info()
				if err != nil {
					errc <- err
					continue
				}
				regexPattern, globPattern := fileExtFilterForLang(info.Language)

				for _, d := range m.dirs {
					files, err := filepath.Glob(path.Join(d, globPattern))
					if err != nil {
						// Non-fatal
						log.Printf("Error reading files in %s: %v\n", d, err)
					}
					for _, f := range files {
						pendingChannel <- pendingReview{configHash, f}
					}
				}

				for _, f := range m.files {
					if matches, err := regexp.MatchString(regexPattern, f); !matches {
						if err != nil {
							log.Println("error in regex: ", regexPattern)
						}
						continue
					}
					pendingChannel <- pendingReview{configHash, f}
				}
			}
		}
		close(pendingChannel)
		close(reviewChannel)
	}()
	return reviewChannel, pendingChannel
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

	// Explode cfg/dirs into service->filepath pairs
	rc, pc := reviewQueue(ctx, configDirs, errc)

	reviews := make(map[string]reviewStream)

	var issueWG sync.WaitGroup
	issueMutex := &sync.Mutex{}

	issues := []*api.Issue{}

l:
	for {
		select {
		case r, open := <-rc:
			if !open {
				continue
			}
			issueWG.Add(1)
			reviews[r.configHash] = r
			go func(r reviewStream) {
				for i := range r.issuesc {
					issueMutex.Lock()
					if cfm.Confirm(0, i) {
						issues = append(issues, i)
					}
					issueMutex.Unlock()
				}
				r.service.Close()
				issueWG.Done()
			}(r)
		case p, open := <-pc:
			if !open {
				for _, r := range reviews {
					close(r.filesc)
				}
				break l
			}
			reviews[p.configHash].filesc <- p.filePath
		}
	}

	issueWG.Wait()

	close(errc)

	outputFmt := review.OutputFormat(ctx.String("output-fmt"))

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

	// Print formatted output, even if there are no issues (eg empty json {})
	if outputFmt != "none" {
		output := review.Output(outputFmt, ctx.String("output"), issues)
		fmt.Print(output)
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
