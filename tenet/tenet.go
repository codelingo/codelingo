package tenet

import (
	"fmt"

	"github.com/codegangsta/cli"

	devTenet "github.com/lingo-reviews/dev/tenet"

	"github.com/lingo-reviews/lingo/tenet/driver"
)

// TODO: Better way to re-export, or way to avoid?
type Config driver.Common

type Tenet interface {
	String() string
	InitDriver() error // TODO: Can this be private?
	Review(args ...string) (*driver.ReviewResult, error)
	Help(args ...string) (string, error)
	Version() (string, error)
	CommentSet() (*devTenet.CommentSet, error)
	Debug(args ...string) string
	GetOptions() driver.Options
	Pull() error
}

// NewTenet builds and returns a Tenet that is not yet initialised.
func New(ctx *cli.Context, cfg Config) (Tenet, error) {
	switch cfg.Driver {
	case "docker", "": // Default driver
		return driver.NewDocker(ctx, driver.Common(cfg))
	case "binary":
		return driver.NewBinary(ctx, driver.Common(cfg))
	}

	return nil, fmt.Errorf("Unknown driver specified: %q", cfg.Driver)
}
