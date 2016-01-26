package commands

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/codegangsta/cli"
	"github.com/mohae/utilitybelt/deepcopy"

	jt "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type CMDTest struct {
	jt.CleanupSuite
	// jt.FakeHomeSuite
	Context *cli.Context
	stdErr  bytes.Buffer
}

var _ = gc.Suite(&CMDTest{})

// 	c.Context = &cli.Context{
// 			App            *App
// 	Command        Command
// 	flagSet        *flag.FlagSet
// 	globalSet      *flag.FlagSet
// 	setFlags       map[string]bool
// 	globalSetFlags map[string]bool
// 	}
// }

var mockTenetCfg = &config{
	TenetGroups: []TenetGroup{
		{Name: "default",
			Tenets: []TenetConfig{
				{
					Name: "lingoreviews/tenetseed:latest",
				}, {
					Name: "lingoreviews/space_after_forward_slash",
				}, {
					Name: "lingo-reviews/unused_function_args",
				}, {
					Name: "lingo-reviews/license",
					Options: map[string]interface{}{
						"header": "// MIT\n",
					},
				},
			},
		},
	},
}

func testCfg(c *gc.C) (cfgPath string, closer func()) {
	f, err := ioutil.TempFile("", "mockTenetCfg")
	c.Assert(err, jc.ErrorIsNil)
	ctx := mockContext(c, tenetCfgFlg.longArg(), f.Name(), "noop")
	c.Assert(writeConfigFile(ctx, mockTenetCfg), jc.ErrorIsNil)
	return f.Name(), func() {
		os.Remove(f.Name())
		f.Close()
	}
}

func longName(f cli.Flag) string {
	parts := strings.Split(f.String(), ",")
	return strings.TrimLeft(parts[0], "-")
}

func addGlobalOpts(set *flag.FlagSet) {
	for _, flg := range GlobalOptions {
		lName := longName(flg)
		switch f := flg.(type) {
		case cli.BoolFlag:
			set.Bool(lName, false, f.Usage)
		case cli.StringFlag:
			set.String(lName, f.Value, f.Usage)
		}
	}
}

// mockContext is a test helper for testing commands. Flags should only be set
// with their long name.
func mockContext(c *gc.C, args ...string) *cli.Context {
	set := flag.NewFlagSet("test", 0)
	addGlobalOpts(set)

	c.Assert(set.Parse(args), jc.ErrorIsNil)

	ctx := cli.NewContext(cli.NewApp(), set, nil)
	ctx.Command = cli.Command{Name: ctx.Args().First()}
	return ctx
}

func (s *CMDTest) SetUpSuite(c *gc.C) {
	origExiter := exiter
	exiter = func(code int) {
		//noOp func
	}
	stderr = &s.stdErr

	s.AddSuiteCleanup(func(c *gc.C) {
		exiter = origExiter
		stderr = os.Stderr
	})
}

func (s *CMDTest) SetUpTest(c *gc.C) {
	// cleanout err buffer
	s.stdErr = bytes.Buffer{}
}

func (*CMDTest) TestWriteAndReadTenetCfg(c *gc.C) {
	f, err := ioutil.TempFile("", "mockTenetCfg")
	fName := f.Name()
	defer func() {
		os.Remove(fName)
		f.Close()
	}()
	c.Assert(err, jc.ErrorIsNil)
	ctx := mockContext(c, tenetCfgFlg.longArg(), fName)
	c.Assert(writeConfigFile(ctx, mockTenetCfg), jc.ErrorIsNil)

	obtained, err := readConfigFile(fName)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(obtained, gc.DeepEquals, mockTenetCfg)
}

func (*CMDTest) TestDesiredTenetCfgPath(c *gc.C) {
	ctx := mockContext(c, tenetCfgFlg.longArg(), "custom/cfg/path")
	c.Assert(desiredTenetCfgPath(ctx), gc.Equals, "custom/cfg/path")
}

func (*CMDTest) TestTenetCfgPath(c *gc.C) {
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
		ctx := mockContext(c, tenetCfgFlg.longArg(), cfgPath)
		cPath, err := tenetCfgPath(ctx)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(cPath, gc.Matches, expected)
	}
}

func (*CMDTest) TestMockContext(c *gc.C) {
	ctx := mockContext(c, tenetCfgFlg.longArg(), "custom/cfg/path", "add")
	c.Assert(ctx.Args(), gc.DeepEquals, cli.Args{"add"})
	c.Assert(ctx.GlobalString(tenetCfgFlg.long), gc.Equals, "custom/cfg/path")
	c.Assert(ctx.Command.Name, gc.Equals, "add")
}

var complexConfig = &config{
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

func (*CMDTest) TestAllTenets(c *gc.C) {
	// AllTenets Should provide a list of all tenets flattened from groups and without duplicates
	output := complexConfig.AllTenets()

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

func (*CMDTest) TestHasTenetGroup(c *gc.C) {
	// HasTenetGroup should return true only for existing group names, case sensitive
	c.Assert(complexConfig.HasTenetGroup("default"), gc.Equals, true)
	c.Assert(complexConfig.HasTenetGroup("Group1"), gc.Equals, true)
	c.Assert(complexConfig.HasTenetGroup("Group2"), gc.Equals, true)
	c.Assert(complexConfig.HasTenetGroup("invalid"), gc.Equals, false)
	c.Assert(complexConfig.HasTenetGroup("group1"), gc.Equals, false)
}

func (*CMDTest) TestAddTenetGroup(c *gc.C) {
	// AddTenetGroup should add a new empty group to the config if that group doesn't already exist
	testConfig := deepcopy.Iface(complexConfig).(*config)

	// Test duplicate
	testConfig.AddTenetGroup("Group1")
	c.Assert(testConfig, gc.DeepEquals, complexConfig)

	// Test new
	expected := &config{
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

func (*CMDTest) TestRemoveTenetGroup(c *gc.C) {
	// RemoveTenetGroup should remove a group and all contained tenets, and have no effect for missing groups
	testConfig := deepcopy.Iface(complexConfig).(*config)

	// Test nonexistent
	testConfig.RemoveTenetGroup("invalid")
	c.Assert(testConfig, gc.DeepEquals, complexConfig)

	// Test existing
	expected := &config{
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

func (*CMDTest) TestAddTenet(c *gc.C) {
	// AddTenet should add a tenet by config to an existing group, create a group if the requested does not exist, and return error
	// for duplicate names/group combos
	testConfig1 := deepcopy.Iface(complexConfig).(*config)
	testConfig2 := deepcopy.Iface(complexConfig).(*config)
	testConfig3 := deepcopy.Iface(complexConfig).(*config)

	// Add tenet to existing group
	expected1 := &config{
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
	expected2 := &config{
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

func (*CMDTest) TestRemoveTenet(c *gc.C) {
	// RemoveTenet should remove a tenet only from the named group, and return an error if either the group or tenet is not found
	testConfig1 := deepcopy.Iface(complexConfig).(*config)
	testConfig2 := deepcopy.Iface(complexConfig).(*config)

	// Test existing
	expected1 := &config{
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

func (*CMDTest) TestHasTenet(c *gc.C) {
	// Test all existing
	c.Assert(complexConfig.HasTenet("lingoreviews/tenetseed:latest"), gc.Equals, true)
	c.Assert(complexConfig.HasTenet("lingo-reviews/license"), gc.Equals, true)
	c.Assert(complexConfig.HasTenet("lingo-reviews/unused_function_args"), gc.Equals, true)
	c.Assert(complexConfig.HasTenet("lingoreviews/space_after_forward_slash"), gc.Equals, true)
	// Test some missing
	c.Assert(complexConfig.HasTenet("lingoreviews/tenetseed"), gc.Equals, false)
	c.Assert(complexConfig.HasTenet("lingoreviews/license"), gc.Equals, false)
	c.Assert(complexConfig.HasTenet("nonamespace/notenet"), gc.Equals, false)
}
