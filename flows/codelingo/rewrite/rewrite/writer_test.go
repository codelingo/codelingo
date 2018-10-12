package rewrite_test

import (
	"github.com/codelingo/codelingo/flows/codelingo/rewrite/rewrite"
	rewriterpc "github.com/codelingo/codelingo/flows/codelingo/rewrite/rpc"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func (s *cmdSuite) TestWrite(c *gc.C) {
	c.Skip("reason")
	hunks := []*rewriterpc.Hunk{{
		Filename:    "test/mock.go",
		StartOffset: int64(19),
		EndOffset:   int64(23),
		SRC:         "newName",
	}}
	err := rewrite.Write(hunks)

	c.Assert(err, jc.ErrorIsNil)

}
