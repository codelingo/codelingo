package review

import (
	"testing"

	jc "github.com/juju/testing/checkers"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type reviewSuite struct{}

var _ = Suite(&reviewSuite{})

func (s *reviewSuite) TestParseURL(c *C) {
	urlStr := "https://github.com/waigani/codelingo_demo/pull/1"
	_, err := ParsePR(urlStr)
	c.Assert(err, jc.ErrorIsNil)
}
