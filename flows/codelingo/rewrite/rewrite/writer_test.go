package rewrite

import (
	"fmt"
	"strings"

	rewriterpc "github.com/codelingo/codelingo/flows/codelingo/rewrite/rpc"

	flowutil "github.com/codelingo/codelingo/sdk/flow"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func (s *cmdSuite) TestWrite(c *gc.C) {
	c.Skip("reason")
	results := []*flowutil.DecoratedResult{
		{
			Ctx: nil,
			Payload: &rewriterpc.Hunk{
				Filename:    "test/mock.go",
				StartOffset: int32(19),
				EndOffset:   int32(23),
				SRC:         "newName",
			}},
	}

	err := Write(results)
	c.Assert(err, jc.ErrorIsNil)

}

func (s *cmdSuite) TestNewFileSRC(c *gc.C) {

	for _, data := range testData {

		hunk := &rewriterpc.Hunk{
			SRC:              "<NEW CODE>",
			StartOffset:      int32(19),
			EndOffset:        int32(23),
			DecoratorOptions: data.decOpts,
			Filename:         "not_used",
		}

		ctx, err := flowutil.NewCtx(&DecoratorApp.App, strings.Split(hunk.DecoratorOptions, " ")[1:])
		c.Assert(err, jc.ErrorIsNil)

		newCode, err := newFileSRC(ctx, hunk, []byte(oldSRC))
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(string(newCode), gc.Equals, string(data.newSRC))
		fmt.Println("PASS:", data.decOpts)

	}
}

var oldSRC string = `
package test

func main() {

}
`[1:]

var testData = []struct {
	decOpts string
	newSRC  []byte
}{
	{
		decOpts: "rewrite \"<NEW CODE>\"",
		newSRC: []byte(`
package test

func <NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite name",
		newSRC: []byte(`
package test

func <NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --replace name",
		newSRC: []byte(`
package test

func <NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --replace --start-to-end-offset name",
		newSRC: []byte(`
package test

func <NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-to-end-offset name",
		newSRC: []byte(`
package test

func <NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset name",
		newSRC: []byte(`
package test

func <NEW CODE>ain() {

}
`[1:]),
	}, {
		decOpts: "rewrite --line name",
		newSRC: []byte(`
package test

<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset --line name",
		newSRC: []byte(`
package test

<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --end-offset --line name",
		newSRC: []byte(`
package test

<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --start-to-end-offset --prepend name",
		newSRC: []byte(`
package test

func <NEW CODE>main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset --prepend name",
		newSRC: []byte(`
package test

func <NEW CODE>main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --prepend name",
		newSRC: []byte(`
package test

func <NEW CODE>main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --end-offset --prepend name",
		newSRC: []byte(`
package test

func mai<NEW CODE>n() {

}
`[1:]),
	}, {
		decOpts: "rewrite --prepend --line name",
		newSRC: []byte(`
package test

<NEW CODE>
func main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-to-end-offset --prepend --line name",
		newSRC: []byte(`
package test

<NEW CODE>
func main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset --prepend --line name",
		newSRC: []byte(`
package test

<NEW CODE>
func main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --end-offset --prepend --line name",
		newSRC: []byte(`
package test

<NEW CODE>
func main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --append name",
		newSRC: []byte(`
package test

func main<NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-to-end-offset --append name",
		newSRC: []byte(`
package test

func main<NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset --append name",
		newSRC: []byte(`
package test

func m<NEW CODE>ain() {

}
`[1:]),
	}, {
		decOpts: "rewrite --end-offset --append name",
		newSRC: []byte(`
package test

func main<NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --append --line name",
		newSRC: []byte(`
package test

func main() {
<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --start-to-end-offset --append --line name",
		newSRC: []byte(`
package test

func main() {
<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --end-offset --append --line name",
		newSRC: []byte(`
package test

func main() {
<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset --append --line name",
		newSRC: []byte(`
package test

func main() {
<NEW CODE>

}
`[1:]),
	},
}

// TODO(waigani) test replace first line
