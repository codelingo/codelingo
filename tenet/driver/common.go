package driver

import "github.com/codegangsta/cli"

type Options map[string]interface{}

type Common struct {
	// Name of the tenet
	Name string `toml:"name"`

	// Name of the driver in use
	Driver string `toml:"driver"`

	// Tag of the tenet
	Tag string `toml:"tag"`

	// Registry server to pull the image from
	Registry string `toml:"registry"`

	// Config options for tenet
	Options Options `toml:"options"`

	context *cli.Context
}

func (c *Common) String() string {
	if c.Tag != "" {
		return c.Name + ":" + c.Tag
	}
	return c.Name
}

func (c *Common) GetOptions() Options {
	return c.Options
}
