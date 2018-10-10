package parse_test

import (
	"testing"

	"github.com/codelingo/clql/dotlingo"
	"github.com/codelingo/clql/inner"
	"github.com/codelingo/platform/flow/bots/parse"

	jc "github.com/juju/testing/checkers"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	TestingT(t)
}

type suite struct{}

// TODO: give suite mock caller - the test currently looks at github.com/codelingo/codelingo
var _ = Suite(&suite{})

func (s *suite) TestTenetImport(c *C) {
	tenet := `
tenets:
  - import: cacophony/default/defer-in-loop
`[1:]
	importer := &dotlingo.Dotlingo{
		Imports: []*dotlingo.Import{
			&dotlingo.Import{
				ImportType: &dotlingo.Import_Tenet{
					Tenet: &dotlingo.TenetImport{
						Owner:  "cacophony",
						Bundle: "default",
						Name:   "defer-in-loop",
					},
				},
			},
		},
	}

	res, err := parse.Parse("", tenet)
	c.Check(err, jc.ErrorIsNil)

	c.Check(res[0], jc.DeepEquals, importer)
	c.Check(len(res), jc.DeepEquals, 2)
}

func (s *suite) TestBundleImport(c *C) {
	tenet := `
tenets:
  - import: codelingo/go
`[1:]
	importer := &dotlingo.Dotlingo{
		Imports: []*dotlingo.Import{
			&dotlingo.Import{
				ImportType: &dotlingo.Import_Bundle{
					Bundle: &dotlingo.BundleImport{
						Owner: "codelingo",
						Name:  "go",
					},
				},
			},
		},
	}

	res, err := parse.Parse("", tenet)
	c.Check(err, jc.ErrorIsNil)
	// error is nil
	//c.Check(err.Error(), jc.Contains, "parse failed for imported tenet codelingo/go")

	_ = importer
	_ = res
	// c.Check(res[0], jc.DeepEquals, importer)
	// c.Check(len(res), jc.DeepEquals, 9)
}

func (s *suite) TestBundleAndTenetImport(c *C) {
	tenet := `
tenets:
  - import: codelingo/go
  - import: cacophony/default/defer-in-loop
`[1:]

	importer := &dotlingo.Dotlingo{
		Imports: []*dotlingo.Import{
			&dotlingo.Import{
				ImportType: &dotlingo.Import_Bundle{
					Bundle: &dotlingo.BundleImport{
						Owner: "codelingo",
						Name:  "go",
					},
				},
			},
			&dotlingo.Import{
				ImportType: &dotlingo.Import_Tenet{
					Tenet: &dotlingo.TenetImport{
						Owner:  "cacophony",
						Bundle: "default",
						Name:   "defer-in-loop",
					},
				},
			},
		},
	}

	res, err := parse.Parse("", tenet)
	c.Check(err, jc.ErrorIsNil)
	// error is nil
	//c.Check(err.Error(), jc.Contains, "parse failed for imported tenet codelingo/go")

	_ = importer
	_ = res
	// c.Check(res[0], jc.DeepEquals, importer)
	// c.Check(len(res), jc.DeepEquals, 10)
}

func (s *suite) TestImportAndExisting(c *C) {
	tenet := `
tenets:
  - import: cacophony/default/defer-in-loop
  - name: existing-test
    doc: Find all go funcs.
    flows:
      codelingo/review:
        comment: New comment
    query: |
      import codelingo/ast/go

      @ review.issue
      go.func
`[1:]

	importer := &dotlingo.Dotlingo{
		Imports: []*dotlingo.Import{
			&dotlingo.Import{
				ImportType: &dotlingo.Import_Tenet{
					Tenet: &dotlingo.TenetImport{
						Owner:  "cacophony",
						Bundle: "default",
						Name:   "defer-in-loop",
					},
				},
			},
		},
		Tenets: []*dotlingo.Tenet{
			&dotlingo.Tenet{
				Name: "existing-test",
				Doc:  "Find all go funcs.",
				Bots: map[string]*dotlingo.Bot{
					"review": {
						Owner: "codelingo",
						Name:  "review",
						Alias: "review",
						Config: map[string]string{
							"comment": "New comment",
						},
					},
				},
				Query: &inner.Query{
					Lexicons: map[string]*inner.Lexicon{
						"go": {
							Owner:   "codelingo",
							Type:    "ast",
							Name:    "go",
							Version: inner.DefaultVersion,
						},
					},
					FactTree: inner.SetRootFact(&inner.Fact{
						Decorators: []*inner.Decorator{
							{
								Namespace: "review",
								Value:     "issue",
							},
						},
						ID: &inner.FactID{
							Namespace: "go",
							Kind:      "func",
						},
						Arguments: &inner.Arguments{
							Depth: &inner.Depth{
								Min: inner.DefaultMinDepth,
								Max: inner.DefaultMaxDepth,
							},
						},
					}),
				},
			},
		},
	}

	res, err := parse.Parse("", tenet)
	c.Check(err, jc.ErrorIsNil)

	c.Check(res[0], jc.DeepEquals, importer)
	c.Check(len(res), jc.DeepEquals, 2)
}
