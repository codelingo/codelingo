package driver

import (
	"fmt"

	"github.com/codegangsta/cli"

	"github.com/lingo-reviews/dev/api"
)

// DryRun is a tenet driver used to perform no action and simply give a
// debug report of what would happen. Intended for internal development
// use only.
type DryRun struct {
	Common
}

func NewDryRun(ctx *cli.Context, cfg Common) (*DryRun, error) {
	return &DryRun{
		Common: Common{
			Driver:  "dryrun",
			Name:    cfg.Name,
			Options: cfg.Options,
			context: ctx,
		},
	}, nil
}

// Do nothing.
func (d *DryRun) InitDriver() error {
	return nil
}

// Log Pull event.
func (d *DryRun) Pull() error {
	fmt.Printf("Pulling tenet: %s\n", d.Name)
	return nil
}

// Log Review event.
func (d *DryRun) Review(args ...string) (*ReviewResult, error) {
	cmd := api.ReviewCMD()
	err := cmd.Flags.Parse(args)
	if err != nil {
		return nil, err
	}
	startFiles := len(args) - cmd.Flags.NArg()

	for _, arg := range args[startFiles:] {
		fmt.Printf("Reviewing file '%s' with tenet: %s %v\n", arg, d.Name, args[:startFiles])
	}
	return &ReviewResult{}, nil
}

// Log Help event.
func (d *DryRun) Help(args ...string) (string, error) {
	fmt.Printf("Calling Help on tenet '%s' with args: %v\n", d.Name, args)
	return "", nil
}

// Log Description event.
func (d *DryRun) Description() (string, error) {
	fmt.Printf("Requesting description of tenet: %s\n", d.Name)
	return "", nil
}

// Log Version event.
func (d *DryRun) Version() (string, error) {
	fmt.Printf("Requesting version of tenet: %s\n", d.Name)
	return "", nil
}

// Do nothing.
func (d *DryRun) Debug(args ...string) string {
	return ""
}
