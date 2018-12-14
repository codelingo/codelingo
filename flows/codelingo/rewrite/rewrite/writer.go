package rewrite

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/codelingo/codelingo/flows/codelingo/rewrite/rewrite/option"
	rewriterpc "github.com/codelingo/codelingo/flows/codelingo/rewrite/rpc"
	flowutil "github.com/codelingo/codelingo/sdk/flow"
	"github.com/urfave/cli"

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

	seenNewFile := make(map[string]bool)

	for filename, results := range resultMap {

		rootPath, err := flowutil.GitCMD("rev-parse", "--show-toplevel")
		if err != nil {
			return errors.Trace(err)
		}

		fullPath := filepath.Join(strings.TrimSuffix(rootPath, "\n"), filename)
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

			ctx := result.Ctx
			hunk := result.Payload.(*rewriterpc.Hunk)

			if ctx.IsSet("new-file") {

				newFileName := ctx.String("new-file")
				if seenNewFile[newFileName] {
					if err != nil {
						return errors.Errorf("cannot add new file %q more than once", newFileName)
					}
				}

				perm := 0755
				if ctx.IsSet("new-file-perm") {
					perm = ctx.Int("new-file-perm")
				}
				if err := ioutil.WriteFile(filepath.Join(filepath.Dir(fullPath), newFileName), []byte(hunk.SRC), os.FileMode(perm)); err != nil {
					return errors.Trace(err)
				}

				seenNewFile[newFileName] = true
				continue
			}

			fileSRC, _, err = newFileSRC(ctx, hunk, fileSRC)
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

type partitionedFile struct {
	srcBeforeStartOffset func() []byte
	srcAfterStartOffset  func() []byte
	srcBeforeEndOffset   func() []byte
	srcAfterEndOffset    func() []byte

	srcBeforeStartLine func() []byte
	srcAfterStartLine  func() []byte
	srcBeforeEndLine   func() []byte
	srcAfterEndLine    func() []byte

	startLineOffsets func() []int32
	endLineOffsets   func() []int32
}

func splitSRC(hunk *rewriterpc.Hunk, fileSRC []byte) partitionedFile {
	startLineOffsets := lineOffsets(fileSRC, hunk.StartOffset)
	endLineOffsets := lineOffsets(fileSRC, hunk.EndOffset)

	return partitionedFile{
		srcBeforeStartOffset: func() []byte { return []byte(string(fileSRC))[0:hunk.StartOffset] },
		srcAfterStartOffset:  func() []byte { return []byte(string(fileSRC))[hunk.StartOffset+1:] },
		srcBeforeEndOffset:   func() []byte { return []byte(string(fileSRC))[0 : hunk.EndOffset-1] },
		srcAfterEndOffset:    func() []byte { return []byte(string(fileSRC))[hunk.EndOffset:] },

		srcBeforeStartLine: func() []byte { return []byte(string(fileSRC))[0:startLineOffsets[0]] },
		srcAfterStartLine:  func() []byte { return []byte(string(fileSRC))[startLineOffsets[1]+1:] },
		srcBeforeEndLine:   func() []byte { return []byte(string(fileSRC))[0:endLineOffsets[0]] },
		srcAfterEndLine:    func() []byte { return []byte(string(fileSRC))[endLineOffsets[1]+1:] },

		startLineOffsets: func() []int32 { return startLineOffsets },
		endLineOffsets:   func() []int32 { return endLineOffsets },
	}
}

type comment struct {
	content string
	// TODO: comments should span multiple lines, but github doesn't allow that https://github.community/t5/How-to-use-Git-and-GitHub/Feature-request-Multiline-reviews-in-pull-requests/m-p/9850#M3225
	line int
}

func newFileSRC(ctx *cli.Context, hunk *rewriterpc.Hunk, fileSRC []byte) ([]byte, *comment, error) {
	parts := splitSRC(hunk, fileSRC)

	rewrittenFile, err := rewriteFile(ctx, fileSRC, []byte(hunk.SRC), parts, hunk)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	var c *comment
	if hunk.Comment != "" {
		commentedSRC, err := rewriteFile(ctx, fileSRC, []byte(hunk.Comment), parts, hunk)
		if err != nil {
			return nil, nil, errors.Trace(err)
		}

		// Find updated line in new rewrittenFile
		for lineNumber, updatedLine := range rewrittenFile {
			if len(commentedSRC) <= lineNumber {
				return nil, nil, errors.New("reached end of commented file before finding updated line")
			}

			commentedLine := commentedSRC[lineNumber]
			if updatedLine != commentedLine {
				c = &comment{
					content: string(commentedLine),
					line:    lineNumber,
				}
				break
			}
		}
	}
	return rewrittenFile, c, nil
}

func rewriteFile(ctx *cli.Context, inputSRC, newSRC []byte, parts partitionedFile, hunk *rewriterpc.Hunk) ([]byte, error) {
	fileSRC := []byte(string(inputSRC))

	opts, err := option.New(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}
	newLine := append(newSRC, '\n')
	newLineAfter := append([]byte{'\n'}, newSRC...)

	switch {
	case opts.IsReplace() && opts.IsStartToEndOffset() && opts.IsByte():
		// replace between start and end bytes
		fileSRC = append(parts.srcBeforeStartOffset(), append(newSRC, parts.srcAfterEndOffset()...)...)

	case opts.IsReplace() && opts.IsStartOffset() && opts.IsByte():
		// replace only the start byte
		fileSRC = append(parts.srcBeforeStartOffset(), append(newSRC, parts.srcAfterStartOffset()...)...)

	case opts.IsReplace() && opts.IsEndOffset() && opts.IsByte():
		// replace only the end byte
		fileSRC = append(parts.srcBeforeEndOffset(), append(newSRC, parts.srcAfterEndOffset()...)...)

	case opts.IsReplace() && opts.IsStartToEndOffset() && opts.IsLine():
		// o.Do(func() {
		fileSRC = append(parts.srcBeforeStartLine(), append(newSRC, parts.srcAfterEndLine()...)...)
		// })
	case opts.IsReplace() && opts.IsStartOffset() && opts.IsLine():
		fileSRC = append(parts.srcBeforeStartLine(), append(newSRC, parts.srcAfterStartLine()...)...)

	case opts.IsReplace() && opts.IsEndOffset() && opts.IsLine():
		// replace whole line
		fileSRC = append(parts.srcBeforeEndLine(), append(newSRC, parts.srcAfterEndLine()...)...)

	case opts.IsPrepend() && opts.IsStartToEndOffset() && opts.IsByte():
		fallthrough
	case opts.IsPrepend() && opts.IsStartOffset() && opts.IsByte():
		// insert before startoffset
		// TODO: remove reference to hunk
		fileSRC = append(parts.srcBeforeStartOffset(), append(newSRC, fileSRC[hunk.StartOffset:]...)...)
	case opts.IsPrepend() && opts.IsEndOffset() && opts.IsByte():
		// insert before endoffset
		fileSRC = append(parts.srcBeforeEndOffset(), append(newSRC, fileSRC[hunk.EndOffset-1:]...)...)

	case opts.IsPrepend() && opts.IsStartToEndOffset() && opts.IsLine():
		fallthrough
	case opts.IsPrepend() && opts.IsStartOffset() && opts.IsLine():
		// insert on new line above startoffset
		fileSRC = append(parts.srcBeforeStartLine(), append(newLine, fileSRC[parts.startLineOffsets()[0]:]...)...)

	case opts.IsPrepend() && opts.IsEndOffset() && opts.IsLine():
		// insert on new line above endoffset
		fileSRC = append(parts.srcBeforeEndLine(), append(newLine, fileSRC[parts.endLineOffsets()[0]:]...)...)

	case opts.IsAppend() && opts.IsStartToEndOffset() && opts.IsByte():
		fallthrough
	case opts.IsAppend() && opts.IsEndOffset() && opts.IsByte():
		// insert after endoffset
		fileSRC = append(fileSRC[0:hunk.EndOffset], append(newSRC, parts.srcAfterEndOffset()...)...)

	case opts.IsAppend() && opts.IsStartOffset() && opts.IsByte():
		// insert after startoffset
		fileSRC = append(fileSRC[0:hunk.StartOffset+1], append(newSRC, parts.srcAfterStartOffset()...)...)

	case opts.IsAppend() && opts.IsStartToEndOffset() && opts.IsLine():
		fallthrough
	case opts.IsAppend() && opts.IsEndOffset() && opts.IsLine():
		// insert on new line after endoffset
		fileSRC = append(fileSRC[0:parts.endLineOffsets()[1]+1], append(newLineAfter, parts.srcAfterEndLine()...)...)

	case opts.IsAppend() && opts.IsStartOffset() && opts.IsLine():
		// insert on new line after startoffset
		fileSRC = append(fileSRC[0:parts.startLineOffsets()[1]+1], append(newLineAfter, parts.srcAfterStartLine()...)...)
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
