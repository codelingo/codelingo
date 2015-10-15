package driver

import (
	"encoding/json"

	devTenet "github.com/lingo-reviews/dev/tenet"
)

type ReviewResult struct {
	TenetName string
	Issues    []*devTenet.Issue
	Errs      []string
}

func decodeResult(name string, result string) (*ReviewResult, error) {
	reviewResult := &ReviewResult{}
	err := json.Unmarshal([]byte(result), reviewResult)
	reviewResult.TenetName = name
	return reviewResult, err
}
