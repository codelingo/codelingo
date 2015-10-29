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
	Description() (string, error)
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

// Any returns an initialised tenet using any driver that is locally available.
func Any(ctx *cli.Context, name string) (Tenet, error) {
	cfg := driver.Common{
		Name: name,
	}

	// Try drivers in order of failure speed
	if t, err := driver.NewBinary(ctx, cfg); err == nil {
		if t.InitDriver() == nil {
			return t, nil
		}
	}

	if t, err := driver.NewDocker(ctx, cfg); err == nil {
		if t.InitDriver() == nil {
			return t, nil
		}
	}

	return nil, fmt.Errorf("No driver available for %s", name)
}
