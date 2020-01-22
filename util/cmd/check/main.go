package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/juju/errors"

	"github.com/codelingo/codelingo/util/check"
)

func main() {
	var r results
	for _, relPath := range targets(os.Args[1:]) {
		fmt.Printf("-> %s ... ", relPath)

		// Check a (potential) rule directory.
		start := time.Now()
		absPath, err := filepath.Abs(relPath)
		if err == nil {
			err = check.TestRule(absPath)
		}
		elapsed := getElapsed(start)
		r.Add(err)

		// Show outcome inline.
		outcome := "???"
		detail := ""
		cause := errors.Cause(err)
		switch cause {
		case nil:
			outcome = "ok"
		case check.ErrFailed, check.ErrBadRule:
			// This is a hack, and depends on TestRule not tracing/annotating.
			detail = err.(*errors.Err).Underlying().Error()
			outcome = "FAIL"
		case check.ErrNoRule, check.ErrNoTests:
			outcome = "SKIP"
		default:
			detail = errors.ErrorStack(err)
			outcome = "ERROR"
		}
		fmt.Println(outcome, elapsed)
		if detail != "" {
			fmt.Println(detail)
		}
	}

	// Show summary.
	fmt.Println(r)
	if r.Count() != r.Success {
		os.Exit(1)
	}
}

// targets converts os Args into a list of directories to inspect.
func targets(args []string) []string {
	search := false
	var filtered []string
	for _, arg := range args {
		if arg == "--search" {
			search = true
		} else {
			filtered = append(filtered, arg)
		}
	}
	if len(filtered) == 0 {
		filtered = []string{"."}
	}
	if !search {
		return filtered
	}
	return find(filtered)
}

// find returns all directories containing codelingo files under all supplied roots.
func find(roots []string) []string {
	found := map[string]bool{}
	check := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		switch base := filepath.Base(path); base {
		case "codelingo.yaml", "codelingo.yml":
			found[filepath.Dir(path)] = true
			return filepath.SkipDir
		}
		return nil
	}
	for _, root := range roots {
		filepath.Walk(root, check)
	}
	result := make([]string, 0, len(found))
	for path := range found {
		result = append(result, path)
	}
	sort.Strings(result)
	return result
}

// getElapsed formats time usefully for output.
func getElapsed(start time.Time) string {
	return fmt.Sprintf("(%.3fs)", time.Since(start).Seconds())
}

// results collects outcomes for summarizing later.
type results struct {
	Success  int
	Failure  int
	Invalid  int
	Untested int
	Error    int
}

func (r *results) Add(err error) {
	switch cause := errors.Cause(err); cause {
	case nil:
		r.Success++
	case check.ErrFailed:
		r.Failure++
	case check.ErrNoRule, check.ErrBadRule:
		r.Invalid++
	case check.ErrNoTests:
		r.Untested++
	default:
		r.Error++
	}
}

func (r results) Count() int {
	return r.Success + r.Failure + r.Invalid + r.Untested + r.Error
}

func (r results) String() string {
	return fmt.Sprintf(`
Inspected %d rule directories.
%d passed tests.
%d failed tests.
%d were invalid.
%d had no tests.
%d went wrong in some other way.
`, r.Count(), r.Success, r.Failure, r.Invalid, r.Untested, r.Error)
}
