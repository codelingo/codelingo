package query_test

import (
	"testing"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

type lSomeSuite struct {
}

var _ = Suite(&lSomeSuite{
	doIngest:   false,
	buildStore: false,
})

func Test(t *testing.T) {
	TestingT(t)
}

func (s *SomeSuite) SetUpSuite(c *C) {
	var err error
	c.Assert(err, jc.ErrorIsNil)
}

func (s *SomeSuite) SetUpTest(c *C) {} // Issue

func (s *SomeSuite) TearDownTest(c *C) {} // Issue
