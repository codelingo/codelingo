package commands

import (
	//"testing"

	. "gopkg.in/check.v1"
)

// func Test(t *testing.T) {
// 	TestingT(t)
// }

type globalSuite struct{}

var _ = Suite(&globalSuite{})

func (g *globalSuite) TestHelpAlias(c *C) {
	// Test long form
	c.Assert(isHelpAlias([]string{"--help"}), Equals, true)
	c.Assert(isHelpAlias([]string{"--other"}), Equals, false)
	// Test short form
	c.Assert(isHelpAlias([]string{"-h"}), Equals, true)
	c.Assert(isHelpAlias([]string{"-o"}), Equals, false)
	// Test too many args
	c.Assert(isHelpAlias([]string{"--help", "--me"}), Equals, false)
	c.Assert(isHelpAlias([]string{"--please", "--help"}), Equals, false)
	c.Assert(isHelpAlias([]string{"--arg1", "--arg2"}), Equals, false)
	// Test no args
	c.Assert(isHelpAlias([]string{}), Equals, false)
}
