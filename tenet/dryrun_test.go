package tenet

import (
	"testing"

	. "gopkg.in/check.v1"

	"github.com/lingo-reviews/dev/api"
	"github.com/lingo-reviews/lingo/tenet/driver"
)

func Test(t *testing.T) {
	TestingT(t)
}

type dryRunSuite struct {
	tenet   Tenet
	service TenetService
}

var _ = Suite(&dryRunSuite{})

func (s *dryRunSuite) SetUpSuite(c *C) {
	var err error

	s.tenet, err = New(nil, &driver.Base{Driver: "dryrun"})
	c.Assert(err, IsNil)

	s.service, err = s.tenet.Service()
	c.Assert(err, IsNil)
}

func (s *dryRunSuite) TestPull(c *C) {
	var err error

	err = s.tenet.Pull(false)
	c.Assert(err, IsNil)

	err = s.tenet.Pull(true)
	c.Assert(err, IsNil)
}

func (s *dryRunSuite) TestStart(c *C) {
	err := s.service.Start()
	c.Assert(err, IsNil)
}

func (s *dryRunSuite) TestStop(c *C) {
	err := s.service.Stop()
	c.Assert(err, IsNil)
}

func (s *dryRunSuite) TestReview(c *C) {
	err := s.service.Review(make(chan string), make(chan *api.Issue, 5))
	c.Assert(err, IsNil)
}

func (s *dryRunSuite) TestInfo(c *C) {
	info, err := s.service.Info()
	c.Assert(err, IsNil)
	c.Assert(info.Name, Equals, "dryrun")
}

func (s *dryRunSuite) TestLanguage(c *C) {
	lang, err := s.service.Language()
	c.Assert(err, IsNil)
	c.Assert(lang, Equals, "")
}
