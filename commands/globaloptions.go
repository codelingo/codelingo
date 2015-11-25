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
		Name:   tenetCfgFlg.String(),
		Value:  defaultTenetCfgPath,
		Usage:  "path to the toml config file that details the tenets used. Defaults to " + defaultTenetCfgPath + " in current directory",
		EnvVar: "LINGO_TENET_CONFIG_NAME",
	}, cli.StringFlag{
		Name:   outputTypeFlg.String(),
		Value:  "plain-text",
		Usage:  "json, json-pretty, yaml, toml or plain-text. If an output-template is set, it takes precedence",
		EnvVar: "LINGO_OUTPUT_TYPE",
	}, cli.StringFlag{
		Name:   outputFlg.String(),
		Value:  "cli",
		Usage:  "filepath to write output to. By default output will be printed to the CLI",
		EnvVar: "LINGO_OUTPUT",
	}, cli.StringFlag{
		Name:   repoURLFlg.String(),
		Value:  "",
		Usage:  "remote repository URL, if not supplied a local repository will be looked for",
		EnvVar: "LINGO_REPO_URL",
	}, cli.StringFlag{
		Name:   repoPathFlg.String(),
		Value:  ".",
		Usage:  "path to local repository, defaults to current directory",
		EnvVar: "LINGO_REPO_PATH",
	}, cli.StringFlag{
		Name:   outputTemplateFlg.String(),
		Value:  "",
		Usage:  "a template for the output format",
		EnvVar: "LINGO_OUTPUT_TEMPLATE",
	}, cli.StringFlag{
		Name:   lingoHomeFlg.String(),
		Value:  defaultLingoHome(),
		Usage:  "a directory of files needed for Lingo to operate",
		EnvVar: "LINGO_HOME",
	}, cli.BoolFlag{
		Name:   diffFlg.String(),
		Usage:  "only report issues found in unstaged, uncommited work",
		EnvVar: "LINGO_DIFF",
		// TODO(waigani) move dump flag to review cmd flag
	}, cli.BoolFlag{
		Name:   dumpFlg.String(),
		Usage:  "By default, Lingo prompts the user to confirm each issue found. dump skips this phase, dumping out all issues found.",
		EnvVar: "LINGO-DUMP",
	},
}
