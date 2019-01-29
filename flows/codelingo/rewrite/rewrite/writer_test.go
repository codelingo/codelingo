package rewrite

import (
	"fmt"
	"os"
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

// TODO: implement once rewrite fname is implemented.
func (s *cmdSuite) TestRewriteFileName(c *gc.C) {

}

func (s *cmdSuite) TestNewFile(c *gc.C) {

	newFile := "new_test.go"

	ctx, err := flowutil.NewCtx(&DecoratorApp.App, "--new-file", newFile, "--new-file-perm", "0755")
	c.Assert(err, jc.ErrorIsNil)

	results := []*flowutil.DecoratedResult{
		{

			Ctx: ctx,
			Payload: &rewriterpc.Hunk{
				SRC:         "package rewrite_test",
				StartOffset: int32(19),
				EndOffset:   int32(23),
				Filename:    "flows/codelingo/rewrite/rewrite/writer_test.go",
			},
		},
	}

	c.Assert(Write(results), jc.ErrorIsNil)

	_, err = os.Stat(newFile)
	c.Assert(os.IsNotExist(err), jc.IsFalse)
	c.Assert(os.Remove(newFile), jc.ErrorIsNil)
}

func (s *cmdSuite) TestNewFileSRC(c *gc.C) {

	for _, data := range testData {

		hunk := &rewriterpc.Hunk{
			SRC:              "<NEW CODE>",
			StartOffset:      int32(19),
			EndOffset:        int32(23),
			DecoratorOptions: data.decOpts,
			Filename:         "not_used",
			Comment:          "<ALT CODE>",
		}

		ctx, err := flowutil.NewCtx(&DecoratorApp.App, strings.Split(hunk.DecoratorOptions, " ")[1:]...)
		c.Assert(err, jc.ErrorIsNil)

		newCode, comment, err := newFileSRC(ctx, hunk, []byte(oldSRC))
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(string(newCode), gc.Equals, string(data.newSRC))
		c.Assert(comment, gc.DeepEquals, data.comment)
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
	comment *comment
}{
	{
		decOpts: "rewrite \"<NEW CODE>\"",
		comment: &comment{
			line:     2,
			content:  "func <ALT CODE>() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func <NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite name",
		comment: &comment{
			line:     2,
			content:  "func <ALT CODE>() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func <NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --replace name",
		comment: &comment{
			line:     2,
			content:  "func <ALT CODE>() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func <NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --replace --start-to-end-offset name",
		comment: &comment{
			line:     2,
			content:  "func <ALT CODE>() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func <NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-to-end-offset name",
		comment: &comment{
			line:     2,
			content:  "func <ALT CODE>() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func <NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset name",
		comment: &comment{
			line:     2,
			content:  "func <ALT CODE>ain() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func <NEW CODE>ain() {

}
`[1:]),
	}, {
		decOpts: "rewrite --line name",
		comment: &comment{
			line:     2,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset --line name",
		comment: &comment{
			line:     2,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --end-offset --line name",
		comment: &comment{
			line:     2,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --start-to-end-offset --prepend name",
		comment: &comment{
			line:     2,
			content:  "func <ALT CODE>main() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func <NEW CODE>main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset --prepend name",
		comment: &comment{
			line:     2,
			content:  "func <ALT CODE>main() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func <NEW CODE>main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --prepend name",
		comment: &comment{
			line:     2,
			content:  "func <ALT CODE>main() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func <NEW CODE>main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --end-offset --prepend name",
		comment: &comment{
			line:     2,
			content:  "func mai<ALT CODE>n() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func mai<NEW CODE>n() {

}
`[1:]),
	}, {
		decOpts: "rewrite --prepend --line name",
		comment: &comment{
			line:     2,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

<NEW CODE>
func main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-to-end-offset --prepend --line name",
		comment: &comment{
			line:     2,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

<NEW CODE>
func main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset --prepend --line name",
		comment: &comment{
			line:     2,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

<NEW CODE>
func main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --end-offset --prepend --line name",
		comment: &comment{
			line:     2,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

<NEW CODE>
func main() {

}
`[1:]),
	}, {
		decOpts: "rewrite --append name",
		comment: &comment{
			line:     2,
			content:  "func main<ALT CODE>() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func main<NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-to-end-offset --append name",
		comment: &comment{
			line:     2,
			content:  "func main<ALT CODE>() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func main<NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset --append name",
		comment: &comment{
			line:     2,
			content:  "func m<ALT CODE>ain() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func m<NEW CODE>ain() {

}
`[1:]),
	}, {
		decOpts: "rewrite --end-offset --append name",
		comment: &comment{
			line:     2,
			content:  "func main<ALT CODE>() {",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func main<NEW CODE>() {

}
`[1:]),
	}, {
		decOpts: "rewrite --append --line name",
		comment: &comment{
			line:     3,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func main() {
<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --start-to-end-offset --append --line name",
		comment: &comment{
			line:     3,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func main() {
<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --end-offset --append --line name",
		comment: &comment{
			line:     3,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func main() {
<NEW CODE>

}
`[1:]),
	}, {
		decOpts: "rewrite --start-offset --append --line name",
		comment: &comment{
			line:     3,
			content:  "<ALT CODE>",
			original: "func main() {",
		},
		newSRC: []byte(`
package test

func main() {
<NEW CODE>

}
`[1:]),
	},
}

// TODO(waigani) test replace first line
