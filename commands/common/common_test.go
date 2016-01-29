package common

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/codegangsta/cli"

	jt "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type commonSuite struct {
	jt.CleanupSuite
	// jt.FakeHomeSuite
	Context *cli.Context
	stdErr  bytes.Buffer
}

var _ = gc.Suite(&commonSuite{})

func (s *commonSuite) SetUpSuite(c *gc.C) {
	origExiter := Exiter
	Exiter = func(code int) {
		//noOp func
	}
	Stderr = &s.stdErr

	s.AddSuiteCleanup(func(c *gc.C) {
		Exiter = origExiter
		Stderr = os.Stderr
	})
}

func (s *commonSuite) SetUpTest(c *gc.C) {
	// cleanout err buffer
	s.stdErr = bytes.Buffer{}
}

func (*commonSuite) TestWriteAndReadTenetCfg(c *gc.C) {
	f, err := ioutil.TempFile("", "MockTenetCfg()")
	fName := f.Name()
	defer func() {
		os.Remove(fName)
		f.Close()
	}()
	c.Assert(err, jc.ErrorIsNil)
	ctx := MockContext(c, TenetCfgFlg.LongArg(), fName)
	c.Assert(WriteConfigFile(ctx, MockTenetCfg()), jc.ErrorIsNil)

	obtained, err := ReadConfigFile(fName)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(obtained, gc.DeepEquals, MockTenetCfg())
}

func (*commonSuite) TestDesiredTenetCfgPath(c *gc.C) {
	ctx := MockContext(c, TenetCfgFlg.LongArg(), "custom/cfg/path")
	c.Assert(DesiredTenetCfgPath(ctx), gc.Equals, "custom/cfg/path")
}

func (*commonSuite) TestTenetCfgPath(c *gc.C) {
	// TODO(waigani) Do what skip says. init .lingo in tmp dir, set default to it.
	c.Skip("Errors if default .lingo cannot be found.")
	defaultPathRegex := `^/home/.*\.lingo_home/\.lingo`
	for cfgPath, expected := range map[string]string{
		// the following cfg files are not found, so the default is returned.
		"rand/path/no/file/.lingo": defaultPathRegex,
		"./.lingo":                 defaultPathRegex,
		".lingo":                   defaultPathRegex,
		"/.lingo":                  defaultPathRegex,
	} {
		ctx := MockContext(c, TenetCfgFlg.LongArg(), cfgPath)
		cPath, err := TenetCfgPath(ctx)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(cPath, gc.Matches, expected)
	}
}

func (*commonSuite) TestMockContext(c *gc.C) {
	ctx := MockContext(c, TenetCfgFlg.LongArg(), "custom/cfg/path", "add")
	c.Assert(ctx.Args(), gc.DeepEquals, cli.Args{"add"})
	c.Assert(ctx.GlobalString(TenetCfgFlg.Long), gc.Equals, "custom/cfg/path")
	c.Assert(ctx.Command.Name, gc.Equals, "add")
}

func complexConfig() *Config {
	return &Config{
		Cascade:  true,
		Include:  "*",
		Template: "template.md",
		TenetGroups: []TenetGroup{
			{
				Name: "default",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/tenetseed:latest",
						Options: map[string]interface{}{"opt1": true, "opt2": "an option"},
					},
				},
			},
			{
				Name:     "Group1",
				Template: "g1template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingo-reviews/license",
						Options: map[string]interface{}{"header": "// Copyright 2016\n"},
					},
					{
						// Note this is in Group1 and Group2
						Name:    "lingo-reviews/unused_function_args",
						Options: map[string]interface{}{},
					},
				},
			},
			{
				Name:     "Group2",
				Template: "g2template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/space_after_forward_slash",
						Options: map[string]interface{}{},
					},
					{
						// Note this is in Group1 and Group2
						Name:    "lingo-reviews/unused_function_args",
						Options: map[string]interface{}{},
					},
				},
			},
		},
	}
}

func (*commonSuite) TestAllTenets(c *gc.C) {
	// AllTenets Should provide a list of all tenets flattened from groups and without duplicates
	output := complexConfig().AllTenets()

	expected := []TenetConfig{
		{
			Name:    "lingoreviews/tenetseed:latest",
			Options: map[string]interface{}{"opt1": true, "opt2": "an option"},
		},
		{
			Name:    "lingo-reviews/license",
			Options: map[string]interface{}{"header": "// Copyright 2016\n"},
		},
		{
			Name:    "lingo-reviews/unused_function_args",
			Options: map[string]interface{}{},
		},
		{
			Name:    "lingoreviews/space_after_forward_slash",
			Options: map[string]interface{}{},
		},
	}

	c.Assert(output, gc.DeepEquals, expected)
}

func (*commonSuite) TestHasTenetGroup(c *gc.C) {
	// HasTenetGroup should return true only for existing group names, case sensitive
	testConfig := complexConfig()
	c.Assert(testConfig.HasTenetGroup("default"), gc.Equals, true)
	c.Assert(testConfig.HasTenetGroup("Group1"), gc.Equals, true)
	c.Assert(testConfig.HasTenetGroup("Group2"), gc.Equals, true)
	c.Assert(testConfig.HasTenetGroup("invalid"), gc.Equals, false)
	c.Assert(testConfig.HasTenetGroup("group1"), gc.Equals, false)
}

func (*commonSuite) TestAddTenetGroup(c *gc.C) {
	// AddTenetGroup should add a new empty group to the config if that group doesn't already exist
	testConfig := complexConfig()

	// Test duplicate
	testConfig.AddTenetGroup("Group1")
	c.Assert(testConfig, gc.DeepEquals, complexConfig())

	// Test new
	expected := &Config{
		Cascade:  true,
		Include:  "*",
		Template: "template.md",
		TenetGroups: []TenetGroup{
			{
				Name: "default",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/tenetseed:latest",
						Options: map[string]interface{}{"opt1": true, "opt2": "an option"},
					},
				},
			},
			{
				Name:     "Group1",
				Template: "g1template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingo-reviews/license",
						Options: map[string]interface{}{"header": "// Copyright 2016\n"},
					},
					{
						// Note this is in Group1 and Group2
						Name:    "lingo-reviews/unused_function_args",
						Options: map[string]interface{}{},
					},
				},
			},
			{
				Name:     "Group2",
				Template: "g2template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/space_after_forward_slash",
						Options: map[string]interface{}{},
					},
					{
						// Note this is in Group1 and Group2
						Name:    "lingo-reviews/unused_function_args",
						Options: map[string]interface{}{},
					},
				},
			},
			{
				Name:     "NewGroup",
				Template: "",
				Tenets:   nil,
			},
		},
	}

	testConfig.AddTenetGroup("NewGroup")
	c.Assert(testConfig, gc.DeepEquals, expected)
}

func (*commonSuite) TestRemoveTenetGroup(c *gc.C) {
	// RemoveTenetGroup should remove a group and all contained tenets, and have no effect for missing groups
	testConfig := complexConfig()

	// Test nonexistent
	testConfig.RemoveTenetGroup("invalid")
	c.Assert(testConfig, gc.DeepEquals, complexConfig())

	// Test existing
	expected := &Config{
		Cascade:  true,
		Include:  "*",
		Template: "template.md",
		TenetGroups: []TenetGroup{
			{
				Name: "default",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/tenetseed:latest",
						Options: map[string]interface{}{"opt1": true, "opt2": "an option"},
					},
				},
			},
			{
				Name:     "Group2",
				Template: "g2template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/space_after_forward_slash",
						Options: map[string]interface{}{},
					},
					{
						// Note this is in Group1 and Group2
						Name:    "lingo-reviews/unused_function_args",
						Options: map[string]interface{}{},
					},
				},
			},
		},
	}

	testConfig.RemoveTenetGroup("Group1")
	c.Assert(testConfig, gc.DeepEquals, expected)
}

func (*commonSuite) TestAddTenet(c *gc.C) {
	// AddTenet should add a tenet by config to an existing group, create a group if the requested does not exist, and return error
	// for duplicate names/group combos
	testConfig1 := complexConfig()
	testConfig2 := complexConfig()
	testConfig3 := complexConfig()

	// Add tenet to existing group
	expected1 := &Config{
		Cascade:  true,
		Include:  "*",
		Template: "template.md",
		TenetGroups: []TenetGroup{
			{
				Name: "default",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/tenetseed:latest",
						Options: map[string]interface{}{"opt1": true, "opt2": "an option"},
					},
				},
			},
			{
				Name:     "Group1",
				Template: "g1template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingo-reviews/license",
						Options: map[string]interface{}{"header": "// Copyright 2016\n"},
					},
					{
						// Note this is in Group1 and Group2
						Name:    "lingo-reviews/unused_function_args",
						Options: map[string]interface{}{},
					},
					{
						Name:    "lingoreviews/new_in_group1",
						Options: map[string]interface{}{"a": true, "b": false},
					},
				},
			},
			{
				Name:     "Group2",
				Template: "g2template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/space_after_forward_slash",
						Options: map[string]interface{}{},
					},
					{
						// Note this is in Group1 and Group2
						Name:    "lingo-reviews/unused_function_args",
						Options: map[string]interface{}{},
					},
				},
			},
		},
	}

	err := testConfig1.AddTenet(TenetConfig{
		Name:    "lingoreviews/new_in_group1",
		Options: map[string]interface{}{"a": true, "b": false},
	}, "Group1")
	c.Assert(testConfig1, gc.DeepEquals, expected1)
	c.Assert(err, gc.IsNil)

	// Add tenet to new group
	expected2 := &Config{
		Cascade:  true,
		Include:  "*",
		Template: "template.md",
		TenetGroups: []TenetGroup{
			{
				Name: "default",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/tenetseed:latest",
						Options: map[string]interface{}{"opt1": true, "opt2": "an option"},
					},
				},
			},
			{
				Name:     "Group1",
				Template: "g1template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingo-reviews/license",
						Options: map[string]interface{}{"header": "// Copyright 2016\n"},
					},
					{
						// Note this is in Group1 and Group2
						Name:    "lingo-reviews/unused_function_args",
						Options: map[string]interface{}{},
					},
				},
			},
			{
				Name:     "Group2",
				Template: "g2template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/space_after_forward_slash",
						Options: map[string]interface{}{},
					},
					{
						// Note this is in Group1 and Group2
						Name:    "lingo-reviews/unused_function_args",
						Options: map[string]interface{}{},
					},
				},
			},
			{
				Name:     "NewGroup",
				Template: "",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/new_in_newgroup",
						Options: nil,
					},
				},
			},
		},
	}

	err = testConfig2.AddTenet(TenetConfig{
		Name:    "lingoreviews/new_in_newgroup",
		Options: nil,
	}, "NewGroup")
	c.Assert(testConfig2, gc.DeepEquals, expected2)
	c.Assert(err, gc.IsNil)

	// Add a duplicated name/group
	err = testConfig3.AddTenet(TenetConfig{
		Name:    "lingoreviews/space_after_forward_slash",
		Options: map[string]interface{}{},
	}, "Group2")
	c.Assert(err, gc.NotNil)
	err = testConfig3.AddTenet(TenetConfig{
		Name:    "lingo-reviews/license",
		Options: nil, // Test different options
	}, "Group1")
	c.Assert(err, gc.NotNil)
}

func (*commonSuite) TestRemoveTenet(c *gc.C) {
	// RemoveTenet should remove a tenet only from the named group, and return an error if either the group or tenet is not found
	testConfig1 := complexConfig()
	testConfig2 := complexConfig()

	// Test existing
	expected1 := &Config{
		Cascade:  true,
		Include:  "*",
		Template: "template.md",
		TenetGroups: []TenetGroup{
			{
				Name: "default",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/tenetseed:latest",
						Options: map[string]interface{}{"opt1": true, "opt2": "an option"},
					},
				},
			},
			{
				Name:     "Group1",
				Template: "g1template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingo-reviews/license",
						Options: map[string]interface{}{"header": "// Copyright 2016\n"},
					},
					{
						// Note this is in Group1 and Group2
						Name:    "lingo-reviews/unused_function_args",
						Options: map[string]interface{}{},
					},
				},
			},
			{
				Name:     "Group2",
				Template: "g2template.md",
				Tenets: []TenetConfig{
					{
						Name:    "lingoreviews/space_after_forward_slash",
						Options: map[string]interface{}{},
					},
				},
			},
		},
	}

	err := testConfig1.RemoveTenet("lingo-reviews/unused_function_args", "Group2")
	c.Assert(testConfig1, gc.DeepEquals, expected1)
	c.Assert(err, gc.IsNil)

	// Test non-existing group
	err = testConfig2.RemoveTenet("lingo-reviews/license", "nogroup")
	c.Assert(err, gc.NotNil)

	// Test non-existing tenet
	err = testConfig2.RemoveTenet("lingo-reviews/notenet", "Group1")
	c.Assert(err, gc.NotNil)
}

func (*commonSuite) TestHasTenet(c *gc.C) {
	testConfig := complexConfig()
	// Test all existing
	c.Assert(testConfig.HasTenet("lingoreviews/tenetseed:latest"), gc.Equals, true)
	c.Assert(testConfig.HasTenet("lingo-reviews/license"), gc.Equals, true)
	c.Assert(testConfig.HasTenet("lingo-reviews/unused_function_args"), gc.Equals, true)
	c.Assert(testConfig.HasTenet("lingoreviews/space_after_forward_slash"), gc.Equals, true)
	// Test some missing
	c.Assert(testConfig.HasTenet("lingoreviews/tenetseed"), gc.Equals, false)
	c.Assert(testConfig.HasTenet("lingoreviews/license"), gc.Equals, false)
	c.Assert(testConfig.HasTenet("nonamespace/notenet"), gc.Equals, false)
}
