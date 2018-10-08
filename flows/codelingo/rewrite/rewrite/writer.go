package rewrite

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/juju/errors"
)

func Write(newSRCs []*SRCHunk) error {

	// TODO(waigani) use one open file handler per file to write all changes
	// and use a buffered writer: https://www.devdungeon.com/content/working-
	// files-go#write_buffered

	// first group all issues by file
	issueMap := make(map[string][]*SRCHunk)

	for _, newSRC := range newSRCs {
		issueMap[newSRC.Filename] = append(issueMap[newSRC.Filename], newSRC)
	}

	for filename, issues := range issueMap {

		fileSRC, err := ioutil.ReadFile(filename)
		if err != nil {
			return errors.Trace(err)
		}

		// then order issues by start offset such that we apply the
		// modifications to the file from the bottom up.
		sort.Sort(byOffset(issues))
		var i int
		var issue *SRCHunk
		for i, issue = range issues {
			fileSRC = append(fileSRC[0:issue.StartOffset], append([]byte(issue.SRC), fileSRC[issue.EndOffset:]...)...)
		}

		if err := ioutil.WriteFile(filename, []byte(fileSRC), 0644); err != nil {
			return errors.Trace(err)
		}
		fmt.Printf("%d modifications made to file %s\n", i, filename)

	}

	return nil
}

type byOffset []*SRCHunk

func (o byOffset) Len() int {
	return len(o)
}

func (o byOffset) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o byOffset) Less(i, j int) bool {
	return o[j].StartOffset < o[i].StartOffset
}
