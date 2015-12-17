// Integration test for review command, designed to test edge cases end-to-end,
// particularly features that are difficult to test in isolation such as
// cascading, mutli-file contexts, diff reviews and grpc to binaries and docker.

package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type reviewSuite struct{}

var _ = Suite(&reviewSuite{})

// Get the real path of the checked out lingo source directory.
func lingoPath(c *C) string {
	_, testFile, _, _ := runtime.Caller(1) // Required because `go test` runs out of a temp dir
	testDir, err := filepath.Abs(filepath.Dir(testFile))
	c.Assert(err, IsNil)
	return filepath.Dir(testDir)
}

// Build and install current lingo.
func installLingo(c *C) {
	cmd := exec.Command("go", "install")
	cmd.Dir = lingoPath(c)
	err := cmd.Run()
	c.Assert(err, IsNil)
}

// Rebuild binary and docker tenets required for testing. We assume that the
// tenets repo is avaialable at lingo-repo/../tenets.
func buildTenets(c *C) {
	basePath := filepath.Join(filepath.Dir(lingoPath(c)), "tenets")

	var requiredTenets = []struct {
		driver string
		path   string
	}{
		{"binary", "go/tenets/simpleseed"},
		{"binary", "go/tenets/license"},
		{"binary", "go/tenets/test/fif_line"},
		{"binary", "go/tenets/test/fif_node"},
		{"docker", "go/tenets/license"},
	}

	for _, r := range requiredTenets {
		fmt.Println("Building:", r.driver, "-", r.path)

		cmd := exec.Command("lingo", "build", r.driver)
		cmd.Dir = filepath.Join(basePath, r.path)
		err := cmd.Run()
		c.Assert(err, IsNil)
	}
}

func (s *reviewSuite) SetUpSuite(c *C) {
	//installLingo(c)
	//buildTenets(c)
}

func (s *reviewSuite) TestAll(c *C) {
	fmt.Println("Running: whole project review")

	cmd := exec.Command("lingo", "review", "--keep-all", "--output-fmt", "plain-text")
	cmd.Dir = "project"
	out, err := cmd.Output()
	c.Assert(err, IsNil)

	expected, err := ioutil.ReadFile("project/expected.txt")
	c.Assert(err, IsNil)

	// Reported issues are order independent - sort before comparing
	o := regexp.MustCompile(" *\\d+. ").Split(string(out), -1)
	sort.Strings(o)
	e := regexp.MustCompile(" *\\d+. ").Split(string(expected), -1)
	sort.Strings(e)

	c.Assert(o, DeepEquals, e)
}

func (s *reviewSuite) TestDiff(c *C) {
	fmt.Println("Running: diff review")

	// Copy new files in place
	err := exec.Command("cp", "project/fif/file1.new", "project/fif/file1.go").Run()
	c.Assert(err, IsNil)
	err = exec.Command("cp", "project/fif/file2.new", "project/fif/file2.go").Run()
	c.Assert(err, IsNil)

	cmd := exec.Command("lingo", "review", "--diff", "--keep-all", "--output-fmt", "plain-text")
	cmd.Dir = "project"
	out, err := cmd.Output()
	c.Assert(err, IsNil)

	expected, err := ioutil.ReadFile("project/expected_diff.txt")
	c.Assert(err, IsNil)

	// Reported issues are order independent - sort before comparing
	o := regexp.MustCompile(" *\\d+. ").Split(string(out), -1)
	sort.Strings(o)
	e := regexp.MustCompile(" *\\d+. ").Split(string(expected), -1)
	sort.Strings(e)

	// Clean up changes
	err = exec.Command("cp", "project/fif/file1.original", "project/fif/file1.go").Run()
	c.Assert(err, IsNil)
	err = exec.Command("cp", "project/fif/file2.original", "project/fif/file2.go").Run()
	c.Assert(err, IsNil)

	c.Assert(o, DeepEquals, e)
}
