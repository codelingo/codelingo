package rewrite

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/briandowns/spinner"
	"github.com/codegangsta/cli"
	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/lingo/service"
	grpcclient "github.com/codelingo/lingo/service/grpc"
	"github.com/codelingo/rpc/flow"
	"github.com/codelingo/rpc/flow/client"
	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

func RequestReview(ctx context.Context, req *flow.ReviewRequest) (chan *flow.Issue, chan error, error) {

	// // mock request result while testing
	// issc := make(chan *flow.Issue)
	// errc := make(chan error)
	// go func() {
	// 	defer close(issc)
	// 	defer close(errc)
	// 	issc <- &flow.Issue{
	// 		CtxBefore: "before",
	// 		CtxAfter:  "after",
	// 		LineText:  "line text",
	// 		Comment:   "this is comment",
	// 		Position: &flow.IssueRange{
	// 			Start: &flow.Position{
	// 				Filename: "main.go",
	// 				Offset:   10,
	// 				Line:     2,
	// 			},
	// 			End: &flow.Position{
	// 				Filename: "main.go",
	// 				Offset:   30,
	// 				Line:     3,
	// 			},
	// 		},
	// 	}
	// }()
	// return issc, errc, nil

	conn, err := service.GrpcConnection(service.LocalClient, service.FlowServer)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	c := client.NewFlowClient(conn)

	// Create context with metadata
	ctx, err = grpcclient.AddUsernameToCtx(ctx)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	payload, err := ptypes.MarshalAny(req)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	// TODO(waigani) refactor review, search and rewrite code out of client completely.
	replyc, runErrc, err := c.Run(ctx, &flow.Request{Flow: "rewrite", Payload: payload})
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	issuec := make(chan *flow.Issue)
	errc := make(chan error)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for err := range runErrc {
			errc <- err
		}
		wg.Done()
	}()

	go func() {
		for reply := range replyc {
			if reply.Error != "" {
				errc <- errors.New(reply.Error)
				continue
			}

			issue := &flow.Issue{}
			err := ptypes.UnmarshalAny(reply.Payload, issue)
			if err != nil {
				errc <- err
				continue
			}

			issuec <- issue
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(issuec)
		close(errc)
	}()

	return issuec, errc, nil
}

func MakeReport(issues []*flow.Issue, format, outputFile string) (string, error) {
	var data []byte
	var err error
	switch format {
	case "json":
		data, err = json.Marshal(issues)
		if err != nil {
			return "", errors.Trace(err)
		}
	case "json-pretty":
		data, err = json.MarshalIndent(issues, " ", " ")
		if err != nil {
			return "", errors.Trace(err)
		}
	default:
		return "", errors.Errorf("Unknown format %q", format)
	}

	if outputFile != "" {
		err = ioutil.WriteFile(outputFile, data, 0775)
		if err != nil {
			return "", errors.Annotate(err, "Error writing issues to file")
		}
		return fmt.Sprintf("Done! %d issues written to %s \n", len(issues), outputFile), nil
	}

	return string(data), nil
}

// Read a codelingo.yaml file from a filepath argument
func ReadDotLingo(ctx *cli.Context) (string, error) {
	var dotlingo []byte

	if filename := ctx.String(util.LingoFile.Long); filename != "" {
		var err error
		dotlingo, err = ioutil.ReadFile(filename)
		if err != nil {
			return "", errors.Trace(err)
		}
	}
	return string(dotlingo), nil
}

func ConfirmIssues(cancel context.CancelFunc, issuec chan *flow.Issue, errorc chan error, keepAll bool, saveToFile string) ([]*SRCHunk, error) {
	defer util.Logger.Sync()

	var confirmedIssues []*SRCHunk
	spnr := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spnr.Start()
	defer spnr.Stop()

	output := saveToFile == ""
	cfm, err := NewConfirmer(output, keepAll, nil)
	if err != nil {
		cancel()
		return nil, errors.Trace(err)
	}

	// If user is manually confirming reviews, set a long timeout.
	timeout := time.After(time.Hour * 1)
	if keepAll {
		timeout = time.After(time.Minute * 5)
	}

l:
	for {
		select {
		case err, ok := <-errorc:
			if !ok {
				errorc = nil
				break
			}

			// Abort review
			cancel()
			util.Logger.Debugf("Review error: %s", errors.ErrorStack(err))
			return nil, errors.Trace(err)
		case iss, ok := <-issuec:
			if !keepAll {
				spnr.Stop()
			}
			if !ok {
				issuec = nil
				break
			}

			// TODO: remove errors from issues; there's a separate channel for that
			if iss.Err != "" {
				// Abort review
				cancel()
				return nil, errors.New(iss.Err)
			}

			// TODO(waigani) attach Tenet to issue
			// DEMOWARE hardcoding clql to one in root dir
			dLingoStr, err := ioutil.ReadFile("codelingo.yaml")
			if err != nil {
				return nil, errors.Trace(err)
			}
			type out struct {
				Tenets []struct {
					Bots  map[string]map[string]interface{}
					Query string
				}
			}

			var dLingo out
			err = yaml.Unmarshal(dLingoStr, &dLingo)
			if err != nil {
				return nil, errors.Trace(err)
			}

			clqlStr := "import codelingo/ast/go\n" + iss.Comment

			// support one set decorator per query.
			srcs, err := ClqlToSrc(string(clqlStr))
			if err != nil {
				return nil, err
			}

			// TODO(waigani) align issue with set id. Until then, we can only support one set decorator
			src := srcs["_default"]

			// add full lines for user confirmation
			confirmSRC, err := fullLineSRC(iss, src)
			if err != nil {
				return nil, errors.Trace(err)
			}

			if cfm.Confirm(0, iss, confirmSRC) {

				issuePos := iss.Position
				startPos := issuePos.GetStart()

				// TODO(waigani) offsets need to be based off rewrite decorator, not review.
				confirmedIssues = append(confirmedIssues, &SRCHunk{
					StartOffset: startPos.Offset,
					EndOffset:   issuePos.GetEnd().Offset,
					SRC:         src,
					Filename:    startPos.Filename,
				})
			}

			if !keepAll {
				spnr.Restart()
			}
		case <-timeout:
			cancel()
			return nil, errors.New("timed out waiting for issue")
		}
		if issuec == nil && errorc == nil {
			break l
		}
	}

	// Stop spinner if it hasn't been stopped already
	if keepAll {
		spnr.Stop()
	}
	return confirmedIssues, nil
}

// returns the full lines of the SRC for the issue.
func fullLineSRC(issue *flow.Issue, newSRC string) (string, error) {

	pos := issue.GetPosition()
	startPos := pos.GetStart()
	endPos := pos.GetEnd()

	file, err := os.Open(startPos.Filename)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var byt int64
	var src []byte
	var strtLineByt int64
	var endLineByt int64
	var foundStart bool
	var foundEnd bool
	for scanner.Scan() {
		src = append(src, append(scanner.Bytes(), []byte("\n")...)...)

		startByt := byt
		endByt := startByt + int64(len(scanner.Bytes())+1) // +1 for \n char

		if startPos.Offset >= startByt && startPos.Offset <= endByt {

			// found start line
			strtLineByt = startByt
			foundStart = true
		}

		if endPos.Offset >= startByt && endPos.Offset <= endByt {

			// found end line
			endLineByt = endByt
			foundEnd = true
		}

		if foundStart && foundEnd {

			// xxx.Print(string(src))
			// fmt.Printf("\n[%d:%d]", strtLineByt, startPos.Offset)
			beforeNewSRC := string(src[strtLineByt:startPos.Offset])
			endNewSRC := string(src[endPos.Offset:endLineByt])

			return beforeNewSRC + newSRC + endNewSRC, nil

		}

		byt = endByt
	}

	return "", errors.Trace(scanner.Err())
}

func NewRange(filename string, startLine, endLine int) *flow.IssueRange {
	start := &flow.Position{
		Filename: filename,
		Line:     int64(startLine),
	}

	end := &flow.Position{
		Filename: filename,
		Line:     int64(endLine),
	}

	return &flow.IssueRange{
		Start: start,
		End:   end,
	}
}

type Options struct {
	// TODO(waigani) validate PullRequest
	PullRequest  string
	FilesAndDirs []string
	Diff         bool   // ctx.Bool("diff") TODO(waigani) this should be a sub-command which proxies to git diff
	SaveToFile   string // ctx.String("save")
	KeepAll      bool   // ctx.Bool("keep-all")
	DotLingo     string // ctx.Bool("lingo-file")
	// TODO(waigani) add KeepAllWithTag. Use this for CLAIR autoreviews
	// TODO(waigani) add streaming json output
}
