package main

import (
	"testing"
	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	TestingT(t)
}

type DummyTestSuite struct{}

var _ = Suite(&DummyTestSuite{})



func (s *DummyTestSuite) TestOne(c *C) {


	c.Assert(3, jc.DeepEquals, 3) // Non Issue

	go func(){
		c.Assert(4, jc.DeepEquals, 2) // Issue
	}()


	testAsync(c) // Non Issue
	go testAsync(c) // Issue
}


func testAsync(c *C) {
	c.Assert(4, jc.DeepEquals, 3)  // Issue only if called on a goroutine, requires callgraph.
}
