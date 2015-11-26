package commands

/*import (
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

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/waigani/diffparser"

	"github.com/lingo-reviews/lingo/commands/review"
	"github.com/lingo-reviews/lingo/tenet"
	"github.com/lingo-reviews/lingo/tenet/driver"
	// TODO: Avoid driver import
	"github.com/lingo-reviews/dev/tenet/log"
)*/

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/codegangsta/cli"

	"github.com/lingo-reviews/dev/api"
	"github.com/lingo-reviews/dev/tenet/log"

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

type cfgMap struct {
	cfg  *config
	dirs []string
}

func readCfgs(contextDirs map[string][]string) <-chan cfgMap {
	out := make(chan cfgMap)
	go func() {
		for cfgPath, dirs := range contextDirs {
			cfg, _ := buildConfig(cfgPath, CascadeUp) // TODO: Handle error
			out <- cfgMap{cfg, dirs}
		}
		close(out)
	}()
	return out
}

type pendingReview struct {
	configHash string
	filePath   string
}

type review struct {
	service tenet.TenetService
	filesc  chan string
	issuesc chan *api.Issue
}

var reviews = make(map[string]review)

func reviewQueue(ctx *cli.Context, mappings <-chan cfgMap) <-chan pendingReview {
	out := make(chan pendingReview)
	go func() {
		for m := range mappings {
			// Glob all the files in the associated dirictories for this config, and assign to each tenet by hash
			for _, tc := range m.cfg.AllTenets() {
				// Instantiate a tenet and run service from the config if we haven't seen it before
				configHash := tc.hash()
				r, found := reviews[configHash]
				if !found {
					tn, _ := newTenet(ctx, tc)     // TODO: Handle error
					service, _ := tn.OpenService() // TODO: Handle error
					r = review{
						service: service,
						filesc:  make(chan string, 1),
						issuesc: make(chan *api.Issue),
						// TODO: go func that closes filesc somehow
					}
					service.Review(r.filesc, r.issuesc)
					reviews[configHash] = r
					fmt.Println("Started:", tc.Name) // TODO: Remove
				} else {
					fmt.Println("Found:", tc.Name) // TODO: Remove
				}

				info, _ := r.service.Info() // TODO: Handle error

				for _, d := range m.dirs {
					_, globPattern := fileExtFilterForLang(info.Language)
					files, _ := filepath.Glob(path.Join(d, globPattern)) // TODO: Handle error
					for _, f := range files {
						out <- pendingReview{configHash, f}
					}
				}
			}
		}
		close(out)
	}()
	return out
}

// type result struct {
// 	tenetName string
// 	service   tenet.TenetService
// 	issuesc   chan *api.Issue
// 	err       error
// }

// func runReviews(queue <-chan pendingReview) <-chan result {

// }

func reviewAction(ctx *cli.Context) {
	fmt.Println("Running review")

	// TODO: file args input, as files and dirs

	// Map of project config filenames -> directories they control
	contextDirs := make(map[string][]string)
	filepath.Walk(".", func(relPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			// Ignore folders beginning with '.', except search root
			// TODO: Flag to turn this behaviour off
			if len(relPath) > 1 && info.Name()[0] == '.' {
				return filepath.SkipDir
			}
			// TODO: Faster technique for finding cfgPath taking advantage of Walk's depth-first search
			//       This implementation recurses upwards for each found dir
			cfgPath, _ := tenetCfgPathRecusive(path.Join(relPath, defaultTenetCfgPath))
			contextDirs[cfgPath] = append(contextDirs[cfgPath], relPath)
		}
		return nil
	})

	// Use a channel to read configs with directory mapping
	configDirs := readCfgs(contextDirs)

	// Explode cfg/dirs into service->filepath pairs
	q := reviewQueue(ctx, configDirs)

	/*go func(q <-chan pendingReview) {
		for p := range q {
			r, _ := reviews[p.configHash] // TODO: Handle error
			r.filesc <- p.filePath
		}
	}(q)*/

	//go func(q <-chan pendingReview) {
	for p := range q {
		fmt.Println(p.configHash, p.filePath)
		r, ok := reviews[p.configHash]
		if !ok {
			log.Fatalln("Queued review has no associated service, this should never happen")
		}

		fmt.Println("sending from review", p.filePath)
		r.filesc <- p.filePath
	}
	//}(q)

	fmt.Println("Queue Done")

	//var issueChannels []chan *api.Issue
	for _, r := range reviews {
		// Since q is empty, it's now safe to close filesc's
		// TODO: Move this to go func, make safely concurrent
		close(r.filesc)
		info, _ := r.service.Info()
		fmt.Println("review", info.Name)
		//issueChannels = append(issueChannels, r.issuesc)
	}

	for n := range merge() {
		fmt.Println(n.Name, n.Comment) // Filename, issue comment
	}

	fmt.Println("Done")
}

func merge() <-chan *api.Issue {
	var wg sync.WaitGroup
	out := make(chan *api.Issue)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(r review) {
		for n := range r.issuesc {
			fmt.Println("fowarding issue")
			out <- n
		}
		r.service.Close() // TODO: Handle error
		wg.Done()
	}
	fmt.Println("waiting on", len(reviews))
	wg.Add(len(reviews))
	for _, r := range reviews {
		go output(r)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

/*
// TODO(waigani) this is ~300 line long func! Break it up.
func reviewAction(ctx *cli.Context) {
	var err error

	// Get this first as it might fail, we want to avoid all other work in that case.
	cfm, err := review.NewConfirmer(ctx, diff)
	if err != nil {
		oserrf(err.Error())
		return
	}

	// setup a chan of results.
	resultsc := make(chan *result, totalTenets)

	r := &reviews{
		fileMatches: map[string][]string{},
		ctx:         ctx,
		files:       files,
	}

	for cfg, tenetCfgs := range rQueue {
		for _, tenetConfig := range tenetCfgs {
			go func(tenetConfig TenetConfig) {
				resultsc <- r.startTenetReview(cfg.buildRoot, tenetConfig)
				wg.Done()
			}(tenetConfig)
		}

	}
	log.Println("colating results")
	issues, errs := allResults(ctx, cfm, resultsc)
	log.Println("review done")
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

type reviews struct {
	// A map of all matching files for dir and extension.
	fileMatches map[string][]string
	ctx         *cli.Context
	files       []string
}

func (rs *reviews) startTenetReview(reviewRoot string, tenetCfg TenetConfig) *result {
	r := &result{
		tenetName: tenetCfg.Name,

		// Setting a buffer will allow the tenet to continue
		// to find and load issues while the user confirms.
		issuesc: make(chan *api.Issue, 5),
	}

	var tn tenet.Tenet
	tn, err := newTenet(rs.ctx, tenetCfg)
	if r.addErrOnErr(err) {
		return r
	}
	s, err := tn.OpenService()
	if r.addErrOnErr(err) {
		return r
	}
	r.service = s

	// setup our chans to send / receive from the service.
	filesc := make(chan string)
	if err := s.Review(filesc, r.issuesc); r.addErrOnErr(err) {
		return r
	}

	// Start sending files to review.
	go func(s tenet.TenetService, filesc chan string) {
		info, err := s.Info()
		if r.addErrOnErr(err) {
			return
		}

		// Find files for this tenet.
		regex, glob := fileExtFilterForLang(info.Language)
		globSearch := path.Join(reviewRoot, glob)

		if rs.fileMatches[globSearch] == nil {
			// Use named files if passed in.
			var fNames []string
			if len(rs.files) > 0 {
				// Add only those files that this tenet is interested in.
				for _, fName := range rs.files {
					log.Printf("checking file %q against regex %q", fName, regex)
					if m, err := regexp.MatchString(regex, fName); !m {
						if err != nil {
							// TODO(waigani) log msg here
						}
						continue
					}
					fNames = append(fNames, fName)
				}
			} else {
				log.Println("globing files with", globSearch)
				fNames, err = filepath.Glob(globSearch)
				if err != nil { // Non-fatal.
					// TODO(waigani) log: no files for this tenet to review.
					log.Println("ERROR reading files", err)
				}
				log.Println("globed files")
				log.Printf("#v,", fNames)
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
					rs.fileMatches[globSearch] = append(rs.fileMatches[globSearch], f)
				} else {
					// TODO(waigani) log here.
				}
				file.Close()
			}

			// sort the order of files so the context set by the tenet will be correct.
			sort.Strings(rs.fileMatches[globSearch])
		}

		// Push files to be reviewed by this tenet onto the chan
		// and close the chan once done.
		for _, fName := range rs.fileMatches[globSearch] {
			filesc <- fName
		}
		close(filesc)
		log.Println("closed files")
	}(s, filesc)

	return r
}

type result struct {
	tenetName string
	service   tenet.TenetService
	issuesc   chan *api.Issue
	err       error
}

func (r *result) addErrOnErr(err error) bool {
	if err != nil {
		r.err = err
	}
	return false
}

// TODO(waigani) TECHDEBT if diff is true, we only report the issues found
// within the diff, even though results contains all issues in the target
// file(s). We need to pass the diff  hunk boundaries to the tenets, it is then
// the tenet's responsibility to only analyse those nodes/lines within the
// hunks.

func fanInIssues(wg *sync.WaitGroup, r *result, allIssuesc chan<- *api.Issue) {
	for {
		select {
		case i, ok := <-r.issuesc:
			if !ok && i == nil {
				log.Println("finished fanning in issues for tenet")
				wg.Done()
				log.Println("closing tenet service")
				r.service.Close()
				return
			}
			allIssuesc <- i
		}
	}
}

// allResults returns all the issues all the tenets found.
func allResults(c *cli.Context, cfm *review.IssueConfirmer, resultsc <-chan *result) ([]*api.Issue, []error) {
	allIssuesc := make(chan *api.Issue)
	confirmedIssuesc := make(chan *api.Issue, 1)
	var errs []error

	// First start listening for any issue and prompting for user confirmation
	// as soon as we get one.
	go func(allIssuesc <-chan *api.Issue) {
		for {
			select {
			case issue, ok := <-allIssuesc:
				if !ok && issue == nil {
					close(confirmedIssuesc)
					return
				}

				if cfm.Confirm(0, issue) {
					confirmedIssuesc <- issue
				}
			case <-time.After(20 * time.Second):
				errs = append(errs, errors.New("timed out"))
			}
		}
	}(allIssuesc)

	// Then start fanning in all issues from all tenets into allIssuesc.
	go func(resultsc <-chan *result, allIssuesc chan<- *api.Issue) {
		wg := &sync.WaitGroup{}
		for {
			select {
			case r, ok := <-resultsc:
				if !ok && r == nil {
					// all results are done.
					log.Println("results closed. waiting for all issues to fan in.")
					// wait for all issues to fan in.

					wg.Wait()
					close(allIssuesc)
					log.Println("closing allIssuesc")
					return
				}
				if r.err != nil {
					errs = append(errs, errors.Annotatef(r.err, "tenet %q", r.tenetName))
				}
				wg.Add(1)
				go fanInIssues(wg, r, allIssuesc)
			}
		}
	}(resultsc, allIssuesc)

	// Read the final result off the confirmed issues chan
	var result []*api.Issue
l:
	for {
		select {
		case issue, ok := <-confirmedIssuesc:
			if !ok && issue == nil {
				break l
			}
			result = append(result, issue)
		}
	}

	return result, errs
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
*/
