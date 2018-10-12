package rewrite

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	rewriterpc "github.com/codelingo/codelingo/flows/codelingo/rewrite/rpc"
	flowutil "github.com/codelingo/codelingo/sdk/flow"

	"github.com/juju/errors"
)

func Write(newSRCs []*rewriterpc.Hunk) error {

	// TODO(waigani) use one open file handler per file to write all changes
	// and use a buffered writer: https://www.devdungeon.com/content/working-
	// files-go#write_buffered

	// first group all hunks by file
	hunkMap := make(map[string][]*rewriterpc.Hunk)

	for _, newSRC := range newSRCs {
		hunkMap[newSRC.Filename] = append(hunkMap[newSRC.Filename], newSRC)
	}

	for filename, hunks := range hunkMap {

		rootPath, err := flowutil.GitCMD("root")
		if err != nil {
			return errors.Trace(err)
		}

		fullPath := strings.TrimSuffix(rootPath, "\n") + "/" + filename
		fileSRC, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return errors.Trace(err)
		}

		// then order hunks by start offset such that we apply the
		// modifications to the file from the bottom up.
		sort.Sort(byOffset(hunks))
		var i int
		var hunk *rewriterpc.Hunk
		for i, hunk = range hunks {
			fileSRC = append(fileSRC[0:hunk.StartOffset], append([]byte(hunk.SRC), fileSRC[hunk.EndOffset:]...)...)
		}

		if err := ioutil.WriteFile(fullPath, []byte(fileSRC), 0644); err != nil {
			return errors.Trace(err)
		}
		fmt.Printf("%d modifications made to file %s\n", i+1, fullPath)

	}

	return nil
}

type byOffset []*rewriterpc.Hunk

func (o byOffset) Len() int {
	return len(o)
}

func (o byOffset) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o byOffset) Less(i, j int) bool {
	return o[j].StartOffset < o[i].StartOffset
}
