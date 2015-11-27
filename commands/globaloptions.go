package commands

import (
	"github.com/codegangsta/cli"
	"github.com/lingo-reviews/lingo/util"
)

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
		"start-dir",
		"s",
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
		Name:   repoPathFlg.String(),
		Value:  ".",
		Usage:  "the directory to operate in, defaults to current directory",
		EnvVar: "LINGO_REPO_PATH",
	},

	cli.StringFlag{
		Name:   lingoHomeFlg.String(),
		Value:  util.MustLingoHome(),
		Usage:  "a directory of files needed for Lingo to operate e.g. logs and binary tenets are stored here",
		EnvVar: "LINGO_HOME",
	},

	// TODO(waigani) implement or drop
	// cli.StringFlag{
	// 	Name:   tenetCfgFlg.String(),
	// 	Value:  defaultTenetCfgPath,
	// 	Usage:  "path to a .lingo to use. Defaults to " + defaultTenetCfgPath + " in current directory",
	// 	EnvVar: "LINGO_TENET_CONFIG_NAME",
	// },
	// cli.StringFlag{
	// 	Name:   outputTemplateFlg.String(),
	// 	Value:  "",
	// 	Usage:  "a template for the output format",
	// 	EnvVar: "LINGO_OUTPUT_TEMPLATE",
	// },
	// cli.StringFlag{
	// 	Name:   repoURLFlg.String(),
	// 	Value:  "",
	// 	Usage:  "remote repository URL, if not supplied a local repository will be looked for",
	// 	EnvVar: "LINGO_REPO_URL",
	// },
	// cli.StringFlag{
	// 	Name:   outputFlg.String(),
	// 	Value:  "cli",
	// 	Usage:  "filepath to write output to. By default output will be printed to the CLI",
	// 	EnvVar: "LINGO_OUTPUT",
	// },
}
