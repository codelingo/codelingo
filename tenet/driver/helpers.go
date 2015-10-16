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

// TODO: Use this function for Help/Version/Debug/etc.
// type callFunc func call(method string, result interface{}, args ...string) error

// func makeCall(c callFunc, cmd string, args []string) (string, error) {
//         var response string
//         if err := c(cmd, &response, args...); err != nil {
//                 return "", err
//         }
//         return response, nil
// }

// called as: makeCall(d.call, "Help", args)
