package review

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/golang/protobuf/ptypes"

	"github.com/codegangsta/cli"
	"github.com/codelingo/lingo/app/util"
	"github.com/codelingo/lingo/service"
	grpcclient "github.com/codelingo/lingo/service/grpc"
	"github.com/codelingo/rpc/flow"
	"github.com/codelingo/rpc/flow/client"
	"github.com/golang/protobuf/proto"
	"github.com/juju/errors"
)

func RequestReview(ctx context.Context, req *flow.ReviewRequest) (chan proto.Message, chan error, error) {
	defer util.Logger.Sync()
	util.Logger.Debug("opening connection to flow server ...")
	conn, err := service.GrpcConnection(service.LocalClient, service.FlowServer)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	util.Logger.Debug("...connection to flow server opened")
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

	util.Logger.Debug("sending request to flow server...")
	replyc, runErrc, err := c.Run(ctx, &flow.Request{Flow: "review", Payload: payload})
	if err != nil {
		return nil, nil, errors.Trace(err)
	}
	util.Logger.Debug("...request to flow server sent. Received reply channel.")

	issuec := make(chan proto.Message)
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
			if reply.IsHeartbeat {
				continue
			}
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

type ReportStrt struct {
	Comment  string
	Filename string
	Line     int
	Snippet  string
}

func MakeReport(cliCtx *cli.Context, issues []*ReportStrt) (string, error) {

	format := cliCtx.String("format")
	outputFile := cliCtx.String("output")

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

// TODO(waigani) simplify representation of Issue.
// https://github.com/codelingo/demo/issues/7
// type Issue struct {
// 	apiIssue
// 	TenetName     string `json:"tenetName,omitempty"`
// 	Discard       bool   `json:"discard,omitempty"`
// 	DiscardReason string `json:"discardReason,omitempty"`
// }

// type apiIssue struct {
// 	// The name of the issue.
// 	TenetName     string            `json:"tenetName,omitempty"`
// 	Discard       bool              `json:"discard,omitempty"`
// 	DiscardReason string            `json:"discardReason,omitempty"`
// 	Name          string            `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
// 	Position      *IssueRange       `protobuf:"bytes,2,opt,name=position" json:"position,omitempty"`
// 	Comment       string            `protobuf:"bytes,3,opt,name=comment" json:"comment,omitempty"`
// 	CtxBefore     string            `protobuf:"bytes,4,opt,name=ctxBefore" json:"ctxBefore,omitempty"`
// 	LineText      string            `protobuf:"bytes,5,opt,name=lineText" json:"lineText,omitempty"`
// 	CtxAfter      string            `protobuf:"bytes,6,opt,name=ctxAfter" json:"ctxAfter,omitempty"`
// 	Metrics       map[string]string `protobuf:"bytes,7,rep,name=metrics" json:"metrics,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
// 	Tags          []string          `protobuf:"bytes,8,rep,name=tags" json:"tags,omitempty"`
// 	Link          string            `protobuf:"bytes,9,opt,name=link" json:"link,omitempty"`
// 	NewCode       bool              `protobuf:"varint,10,opt,name=newCode" json:"newCode,omitempty"`
// 	Patch         string            `protobuf:"bytes,11,opt,name=patch" json:"patch,omitempty"`
// 	Err           string            `protobuf:"bytes,12,opt,name=err" json:"err,omitempty"`
// }

// type IssueRange struct {
// 	Start *Position `protobuf:"bytes,1,opt,name=start" json:"start,omitempty"`
// 	End   *Position `protobuf:"bytes,2,opt,name=end" json:"end,omitempty"`
// }

// type Position struct {
// 	Filename string `protobuf:"bytes,1,opt,name=filename" json:"filename,omitempty"`
// 	Offset   int64  `protobuf:"varint,2,opt,name=Offset" json:"Offset,omitempty"`
// 	Line     int64  `protobuf:"varint,3,opt,name=Line" json:"Line,omitempty"`
// 	Column   int64  `protobuf:"varint,4,opt,name=Column" json:"Column,omitempty"`
// }

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
