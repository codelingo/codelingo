package commands

import "github.com/codegangsta/cli"

type flagName struct {
	long  string
	short string
}

func (f *flagName) longArg() string {
	return "--" + f.long
}

func (f *flagName) shortArg() string {
	return "-" + f.short
}

var (

	// global flags
	tenetCfgFlg = flagName{
		"tenet-config",
		"c",
	}
	outputTypeFlg = flagName{
		"output-type",
		"t",
	}
	outputFlg = flagName{
		"output",
		"o",
	}
	repoURLFlg = flagName{
		"repo-url",
		"u",
	}
	repoPathFlg = flagName{
		"repo-path",
		"p",
	}
	outputTemplateFlg = flagName{
		"output-template",
		"e",
	}
	lingoHomeFlg = flagName{
		"lingo-home",
		"l",
	}
	diffFlg = flagName{
		"diff",
		"d",
	}
	dumpFlg = flagName{
		"dump",
		"m",
	}

	//local flags
	allFlg = flagName{
		"all",
		"a",
	}
	updateFlg = flagName{
		"update",
		"u",
	}
	tagsFlg = flagName{
		"tags",
		"g",
	}
	registryFlg = flagName{
		"registry",
		"r",
	}
	driverFlg = flagName{
		"driver",
		"i",
	}
)

func (f *flagName) String() string {
	return f.long + ", " + f.short
}

var GlobalOptions = []cli.Flag{
	cli.StringFlag{
		tenetCfgFlg.String(),
		defaultTenetCfgPath,
		"path to the toml config file that details the tenets used. Defaults to " + defaultTenetCfgPath + " in current directory",
		"LINGO_TENET_CONFIG_NAME",
	}, cli.StringFlag{
		outputTypeFlg.String(),
		"plain-text",
		"json, json-pretty, yaml, toml or plain-text. If an output-template is set, it takes precedence",
		"LINGO_OUTPUT_TYPE",
	}, cli.StringFlag{
		outputFlg.String(),
		"cli",
		"filepath to write output to. By default output will be printed to the CLI",
		"LINGO_OUTPUT",
	}, cli.StringFlag{
		repoURLFlg.String(),
		"",
		"remote repository URL, if not supplied a local repository will be looked for",
		"LINGO_REPO_URL",
	}, cli.StringFlag{
		repoPathFlg.String(),
		".",
		"path to local repository, defaults to current directory",
		"LINGO_REPO_PATH",
	}, cli.StringFlag{
		outputTemplateFlg.String(),
		"",
		"a template for the output format",
		"LINGO_OUTPUT_TEMPLATE",
	}, cli.StringFlag{
		lingoHomeFlg.String(),
		defaultLingoHome(),
		"a directory of files needed for Lingo to operate",
		"LINGO_HOME",
	}, cli.BoolFlag{
		diffFlg.String(),
		"only report issues found in unstaged, uncommited work",
		"LINGO_DIFF",
		// TODO(waigani) move dump flag to review cmd flag
	}, cli.BoolFlag{
		dumpFlg.String(),
		"By default, Lingo prompts the user to confirm each issue found. dump skips this phase, dumping out all issues found.",
		"LINGO-DUMP",
	},
}
