package check

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"

	"github.com/codelingo/clql/dotlingo"
	"github.com/juju/errors"
	"github.com/juju/testing/checkers"
)

var (
	ErrNoRule  = errors.New("not a rule")
	ErrBadRule = errors.New("rule not valid")
	ErrNoTests = errors.New("has no tests")
	ErrFailed  = errors.New("test failed")
)

// TestRule returns an error unless the supplied path contains a valid
// codelingo rule with tests that pass. Specific conditions will be signalled
// by returning errors with a Cause of ErrNoRule, ErrBadRule, ErrNoTests, or
// ErrFailed; no particular inferences can be drawn from other non-nil errors.
func TestRule(dir string) error {
	expect, err := CheckRule(dir)
	if err != nil {
		return err
	}
	actual, err := RunReview(dir)
	if err != nil {
		return err
	}
	if _, err := checkers.DeepEqual(actual, expect); err != nil {
		return errors.Wrap(err, ErrFailed)
	}
	return nil
}

// CheckRule returns the supplied rule directory's expected test results, or an
// error. It will return ErrNoRule if the codelingo file is missing; ErrBadRule
// if the codelingo file cannot be parsed; and ErrNoTests if the expected.json
// file is missing.
func CheckRule(dir string) ([]Issue, error) {
	path := filepath.Join(dir, "codelingo.yaml")
	content, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		path = filepath.Join(dir, "codelingo.yml")
		content, err = ioutil.ReadFile(path)
	}
	if os.IsNotExist(err) {
		return nil, ErrNoRule
	} else if err != nil {
		return nil, errors.Annotatef(err, "cannot read rule file")
	}
	if _, err := dotlingo.Parse(string(content)); err != nil {
		return nil, errors.Wrap(err, ErrBadRule)
	}

	expect, err := ReadIssues(filepath.Join(dir, "expected.json"))
	if os.IsNotExist(errors.Cause(err)) {
		return nil, ErrNoTests
	} else if err != nil {
		return nil, errors.Annotatef(err, "cannot read expected results")
	}
	return expect, nil
}

// RunReview runs `lingo run review` in a temporary git repository copied from dir.
func RunReview(dir string) ([]Issue, error) {
	datePrefix := time.Now().Format("2006-01-02-")
	tempDir, err := ioutil.TempDir("", datePrefix)
	if err != nil {
		return nil, errors.Annotatef(err, "cannot create work dir")
	}
	defer os.RemoveAll(tempDir)

	copyDir := filepath.Join(tempDir, "rule")
	if _, err := Run("", "cp", "-r", dir, copyDir); err != nil {
		return nil, errors.Annotatef(err, "cannot copy rule dir")
	}
	if err := Script(copyDir, [][]string{
		{"git", "init"},
		{"git", "add", "."},
		{"git", "commit", "-m", "for testing"},
		{"lingo", "run", "review", "--keep-all", "-o", "../actual.json"},
	}); err != nil {
		return nil, errors.Annotatef(err, "cannot setup/run lingo review")
	}
	issues, err := ReadIssues(filepath.Join(tempDir, "actual.json"))
	if os.IsNotExist(errors.Cause(err)) {
		// lingo doesn't bother to write output if it found nothing, pretend we
		// found an empty list instead for consistency's sake.
		return nil, nil
	} else if err != nil {
		return nil, errors.Annotatef(err, "cannot read actual results")
	}
	return issues, nil
}

// Run invokes a command in a directory, and returns its CombinedOutput.
func Run(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Annotatef(err, "cannot run; output: %s", out)
	}
	return string(out), nil
}

// Script invokes a sequence of commands, until one fails.
func Script(dir string, cmds [][]string) error {
	for _, cmd := range cmds {
		name, args := cmd[0], cmd[1:]
		if _, err := Run(dir, name, args...); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

// ReadIssues reads a file containing a json-encoded list of Issues, and
// returns them sorted by location.
func ReadIssues(path string) ([]Issue, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Annotatef(err, "cannot read file")
	}
	var issues []Issue
	if err := json.Unmarshal(content, &issues); err != nil {
		return nil, errors.Annotatef(err, "cannot unmarshal content")
	}
	sort.Sort(byLocation(issues))
	return issues, nil
}

// Issue matches the structure of lingo review results.
type Issue struct {
	Comment  string
	Filename string
	Line     int
	Snippet  string
}

// byLocation sorts Issues by file then line.
type byLocation []Issue

func (b byLocation) Len() int      { return len(b) }
func (b byLocation) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b byLocation) Less(i, j int) bool {
	if b[i].Filename < b[j].Filename {
		return true
	}
	return b[i].Line < b[j].Line
}
