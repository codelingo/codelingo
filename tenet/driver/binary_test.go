package driver

import (
	"testing"

	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type binarySuite struct {
	driver Driver
}

var _ = gc.Suite(&binarySuite{})

func (b *binarySuite) SetUpTest(c *gc.C) {
	b.driver = &Binary{Base: &Base{
		Driver: "binary",
		Source: "http://github.com/juju/tenets/juju/workers/nostate",
	},
	}
}

func (b *binarySuite) TestPull(c *gc.C) {
	// Currently real source code will be downloaded and built. This makes the
	// test too brittle.
	c.Skip("need to mock out download and build.")
	err := b.driver.Pull(false)
	c.Assert(err, gc.IsNil)

	err = b.driver.Pull(true)
	c.Assert(err, gc.IsNil)
}
