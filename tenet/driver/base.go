package driver

import (
	"github.com/lingo-reviews/tenets/go/dev/api"
)

type Options map[string]interface{}

type Base struct {
	// Name of the tenet
	Name string

	Driver string

	Registry string

	// TODO(waigani) this is docker specific. Should it be set in
	// ConfigOptions? Or do we have driverOptions?
	Tag string

	// Config options for tenet
	ConfigOptions Options

	// service supports operations on the backing micro-service tenet server.
	service Service
}

func (b *Base) EditFilename(filename string) (editedFilename string) {
	return filename
}

func (b *Base) EditIssue(issue *api.Issue) (editedIssue *api.Issue) {
	return issue
}

func (b *Base) GetOptions() Options {
	return b.ConfigOptions
}

func (b *Base) Pull(bool) error {
	return nil
}
