package tenet

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"

	"github.com/lingo-reviews/lingo/tenet/driver"
	"github.com/lingo-reviews/tenets/go/dev/api"
)

func Test(t *testing.T) {
	TestingT(t)
}

type dryRunSuite struct {
	tenet   Tenet
	service TenetService
}

var _ = Suite(&dryRunSuite{})

func (s *dryRunSuite) SetUpTest(c *C) {
	var err error

	s.tenet, err = New(nil, &driver.Base{Driver: "dryrun"})
	c.Assert(err, IsNil)

	s.service, err = s.tenet.OpenService()
	c.Assert(err, IsNil)
}

func (s *dryRunSuite) TestPull(c *C) {
	err := s.tenet.Pull(false)
	c.Assert(err, IsNil)

	err = s.tenet.Pull(true)
	c.Assert(err, IsNil)
}

func (s *dryRunSuite) TestClose(c *C) {
	err := s.service.Close()
	c.Assert(err, IsNil)
}

func (s *dryRunSuite) TestReview(c *C) {
	filenames := []string{"f1.go", "f2.go", "f3.go"}
	files := make(chan string)
	issues := make(chan *api.Issue, 5)

	go func() {
		err := s.service.Review(files, issues)
		c.Assert(err, IsNil)
	}()

	for _, f := range filenames {
		files <- f
	}

	close(files)

	seenFilenames := []string{}
l:
	for {
		select {
		case issue, ok := <-issues:
			if !ok {
				// issues closed, we're done.
				break l
			}
			// Check all filenames seen and in correct order on completion
			seenFilenames = append(seenFilenames, issue.Position.Start.Filename)
			c.Assert(issue.Name, Equals, "dryrun")
			c.Assert(issue.Comment, Equals, "Dry Run Issue")
			c.Assert(issue.LineText, Equals, "Your code here")
		case <-time.After(3 * time.Second):
			c.Fatal("timed out waiting for issues")
			break l
		}
	}

	c.Assert(filenames, DeepEquals, seenFilenames)
}

func (s *dryRunSuite) TestInfo(c *C) {
	info, err := s.service.Info()
	c.Assert(err, IsNil)
	c.Assert(info, DeepEquals, &api.Info{
		Name:        "dryrun",
		Usage:       "test lingo and configurations",
		Description: "test lingo and configurations ... description",
		Language:    "*",
		Version:     "0.1.0",
	})
}
