package query_test

import (
	"testing"

	"github.com/codelingo/platform/controller/graphdb"
	"github.com/codelingo/platform/tests/setup"
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
	if s.buildStore {
		s.ts, err = setup.NewTestStore()
	} else {
		s.ts, err = setup.NewDgraphTestClient()
	}
	c.Assert(err, jc.ErrorIsNil)

	store, err := graphdb.NewStore()
	c.Assert(err, jc.ErrorIsNil)
	s.store = store
	s.isIngested = make(map[string]bool)
}

func (s *SomeSuite) SetUpTest(c *C) {} // Issue

func (s *SomeSuite) TearDownTest(c *C) {} // Issue
