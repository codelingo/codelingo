package commands

import (
	flg "flag"

	"fmt"
)

type tenetCommand struct {
	Name string
	// A short description of the usage of this command
	Usage string
	// A longer explanation of how the command works
	Description string

	Flags flg.FlagSet // TODO: Can this be private once drivers are in?
}

func tenetCommands() []*tenetCommand {
	return []*tenetCommand{
		reviewCMD(),
		helpCMD(),
		infoCMD(),
		descriptionCMD(),
		versionCMD(),
	}
}

func reviewCMD() *tenetCommand {
	cmd := &tenetCommand{
		Name:  "review",
		Usage: "review target code against this tenet",
	}
	cmd.addFlags(
		Flag{"options", "{}", "json serialized tenet options"},
	)
	return cmd
}

func descriptionCMD() *tenetCommand {
	return &tenetCommand{
		Name:  "description",
		Usage: "gives extended information about this tenet",
	}
}

func helpCMD() *tenetCommand {
	return &tenetCommand{
		Name:  "help",
		Usage: "show help",
	}
}

func infoCMD() *tenetCommand {
	return &tenetCommand{
		Name:  "info",
		Usage: "show information about this tenet",
	}
}

func versionCMD() *tenetCommand {
	return &tenetCommand{
		Name:  "version",
		Usage: "diplays the version of this tenet",
	}
}

type Flag struct {
	Name  string
	Value interface{}
	Usage string
}

func (c *tenetCommand) GetFlag(name string) string {
	f := c.Flags.Lookup(name)
	if f == nil {
		panic(fmt.Sprintf("flag %q not set", name))
	}
	return f.Value.String()
}

func (c *tenetCommand) addFlag(f Flag) *tenetCommand {
	switch v := f.Value.(type) {
	case bool:
		c.Flags.Bool(f.Name, v, f.Usage)
	case string:
		c.Flags.String(f.Name, v, f.Usage)
	default:
		panic(fmt.Sprintf("unrecognised flag value type: %T", v))
	}
	return c
}

func (c *tenetCommand) addFlags(flgs ...Flag) *tenetCommand {
	for _, f := range flgs {
		c.addFlag(f)
	}
	return c
}
