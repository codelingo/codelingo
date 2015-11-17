package driver

import (
	"fmt"

	"github.com/lingo-reviews/lingo/tenet/service"
)

// DryRun is a tenet driver used to perform no action and simply give a
// debug report of what would happen. Intended for internal development
// use only.
type DryRun struct {
	*Base
}

// Log Pull event.
func (d *DryRun) Pull(bool) error {
	fmt.Printf("Pulling tenet: %s\n", d.Name)
	return nil
}

// Init the service.
func (d *DryRun) Service() (service.Service, error) {
	fmt.Printf("Creating dry run service")

	// Matt, I broke your dryRun :( You'll have to create a mock
	// service.Service here.
	return nil, nil
}
