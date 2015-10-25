package tenet

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/lingo-reviews/lingo/tenet/driver"
)

type Config driver.Common

type Tenet interface {
	String() string
	InitDriver() error
	Review(args ...string) (*driver.ReviewResult, error) // TODO: Refactor to not expose driver to callers
	Help(args ...string) (string, error)
	Version() (string, error)
	Debug(args ...string) string
	GetOptions() driver.Options // TODO: Use AddOptions instead to apply cli json args
	Pull() error
}

// NewTenet builds and returns a Tenet that is not yet initialised.
func New(ctx *cli.Context, cfg Config) (Tenet, error) {
	switch cfg.Driver {
	case "docker", "": // Default driver
		return driver.NewDocker(ctx, driver.Common(cfg))
	case "dryrun":
		return driver.NewDryRun(ctx, driver.Common(cfg))
	case "binary":
		return driver.NewBinary(ctx, driver.Common(cfg))
	}

	return nil, fmt.Errorf("Unknown driver specified: %q", cfg.Driver)
}
