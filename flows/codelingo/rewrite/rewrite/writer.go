package rewrite

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/codelingo/codelingo/flows/codelingo/rewrite/rewrite/option"
	rewriterpc "github.com/codelingo/codelingo/flows/codelingo/rewrite/rpc"
	flowutil "github.com/codelingo/codelingo/sdk/flow"

	"github.com/juju/errors"
)

func Write(results []*flowutil.DecoratedResult) error {

	// TODO(waigani) use one open file handler per file to write all changes
	// and use a buffered writer: https://www.devdungeon.com/content/working-
	// files-go#write_buffered

	// first group all results by file
	resultMap := make(map[string][]*flowutil.DecoratedResult)

	for _, result := range results {
		filename := result.Payload.(*rewriterpc.Hunk).Filename
		resultMap[filename] = append(resultMap[filename], result)
	}

	for filename, results := range resultMap {

		rootPath, err := flowutil.GitCMD("rev-parse", "--show-toplevel")
		if err != nil {
			return errors.Trace(err)
		}

		fullPath := strings.TrimSuffix(rootPath, "\n") + "/" + filename
		fileSRC, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return errors.Trace(err)
		}

		// then order results by start offset such that we apply the
		// modifications to the file from the bottom up.
		sort.Sort(byOffset(results))
		var i int
		var result *flowutil.DecoratedResult
		for i, result = range results {

			fileSRC, err = newFileSRC(result.Ctx, result.Payload.(*rewriterpc.Hunk), fileSRC)
			if err != nil {
				return errors.Trace(err)
			}

		}

		if err := ioutil.WriteFile(fullPath, []byte(fileSRC), 0644); err != nil {
			return errors.Trace(err)
		}
		fmt.Printf("%d modifications made to file %s\n", i+1, fullPath)

	}

	return nil
}

// return start and end of the line containing the given offset
func lineOffsets(src []byte, offset int32) []int32 {
	var start, end int32
	// find start
	for i := offset; i >= 0; i-- {
		if src[i] == '\n' {
			break
		}
		start = i
	}

	// find end
	for i := offset; i < int32(len(src)); i++ {
		if src[i] == '\n' {
			break
		}
		end = i
	}
	return []int32{start, end}
}

func newFileSRC(ctx *cli.Context, hunk *rewriterpc.Hunk, fileSRC []byte) ([]byte, error) {

	opts, err := option.New(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}

	newSRC := []byte(hunk.SRC)
	newLine := append(newSRC, '\n')
	newLineAfter := append([]byte{'\n'}, newSRC...)

	startLineOffsets := lineOffsets(fileSRC, hunk.StartOffset)
	endLineOffsets := lineOffsets(fileSRC, hunk.EndOffset)

	srcBeforeStartOffset := fileSRC[0:hunk.StartOffset]
	srcAfterStartOffset := fileSRC[hunk.StartOffset+1:]
	srcBeforeEndOffset := fileSRC[0 : hunk.EndOffset-1]
	srcAfterEndOffset := fileSRC[hunk.EndOffset:]

	srcBeforeStartLine := fileSRC[0:startLineOffsets[0]]
	srcAfterStartLine := fileSRC[startLineOffsets[1]+1:]
	srcBeforeEndLine := fileSRC[0:endLineOffsets[0]]
	srcAfterEndLine := fileSRC[endLineOffsets[1]+1:]

	switch {
	case opts.IsReplace() && opts.IsStartToEndOffset() && opts.IsByte():
		// replace between start and end bytes
		fileSRC = append(srcBeforeStartOffset, append(newSRC, srcAfterEndOffset...)...)

	case opts.IsReplace() && opts.IsStartOffset() && opts.IsByte():
		// replace only the start byte
		fileSRC = append(srcBeforeStartOffset, append(newSRC, fileSRC[hunk.StartOffset+1:]...)...)

	case opts.IsReplace() && opts.IsEndOffset() && opts.IsByte():
		// replace only the end byte
		fileSRC = append(fileSRC[0:hunk.EndOffset-1], append(newSRC, srcAfterEndOffset...)...)

	case opts.IsReplace() && opts.IsStartToEndOffset() && opts.IsLine():
		fileSRC = append(srcBeforeStartLine, append(newSRC, srcAfterEndLine...)...)

	case opts.IsReplace() && opts.IsStartOffset() && opts.IsLine():
		fileSRC = append(srcBeforeStartLine, append(newSRC, srcAfterStartLine...)...)

	case opts.IsReplace() && opts.IsEndOffset() && opts.IsLine():
		// replace whole line
		fileSRC = append(srcBeforeEndLine, append(newSRC, srcAfterEndLine...)...)

	case opts.IsPrepend() && opts.IsStartToEndOffset() && opts.IsByte():
		fallthrough
	case opts.IsPrepend() && opts.IsStartOffset() && opts.IsByte():
		// insert before startoffset
		fileSRC = append(srcBeforeStartOffset, append(newSRC, fileSRC[hunk.StartOffset:]...)...)

	case opts.IsPrepend() && opts.IsEndOffset() && opts.IsByte():
		// insert before endoffset
		fileSRC = append(srcBeforeEndOffset, append(newSRC, fileSRC[hunk.EndOffset-1:]...)...)

	case opts.IsPrepend() && opts.IsStartToEndOffset() && opts.IsLine():
		fallthrough
	case opts.IsPrepend() && opts.IsStartOffset() && opts.IsLine():
		// insert on new line above startoffset
		fileSRC = append(srcBeforeStartLine, append(newLine, fileSRC[startLineOffsets[0]:]...)...)

	case opts.IsPrepend() && opts.IsEndOffset() && opts.IsLine():
		// insert on new line above endoffset
		fileSRC = append(srcBeforeEndLine, append(newLine, fileSRC[endLineOffsets[0]:]...)...)

	case opts.IsAppend() && opts.IsStartToEndOffset() && opts.IsByte():
		fallthrough
	case opts.IsAppend() && opts.IsEndOffset() && opts.IsByte():
		// insert after endoffset
		fileSRC = append(fileSRC[0:hunk.EndOffset], append(newSRC, srcAfterEndOffset...)...)

	case opts.IsAppend() && opts.IsStartOffset() && opts.IsByte():
		// insert after startoffset
		fileSRC = append(fileSRC[0:hunk.StartOffset+1], append(newSRC, srcAfterStartOffset...)...)

	case opts.IsAppend() && opts.IsStartToEndOffset() && opts.IsLine():
		fallthrough
	case opts.IsAppend() && opts.IsEndOffset() && opts.IsLine():
		// insert on new line after endoffset
		fileSRC = append(fileSRC[0:endLineOffsets[1]+1], append(newLineAfter, srcAfterEndLine...)...)

	case opts.IsAppend() && opts.IsStartOffset() && opts.IsLine():
		// insert on new line after startoffset
		fileSRC = append(fileSRC[0:startLineOffsets[1]+1], append(newLineAfter, srcAfterStartLine...)...)
	}

	return fileSRC, nil
}

type byOffset []*flowutil.DecoratedResult

func (o byOffset) Len() int {
	return len(o)
}

func (o byOffset) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o byOffset) Less(i, j int) bool {
	return o[j].Payload.(*rewriterpc.Hunk).StartOffset < o[i].Payload.(*rewriterpc.Hunk).StartOffset
}
