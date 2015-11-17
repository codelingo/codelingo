package driver

import (
	"github.com/lingo-reviews/dev/api"
	"github.com/lingo-reviews/lingo/tenet/service"
)

type Driver interface {
	Pull(bool) error

	Service() (service.Service, error)

	EditFilename(filename string) (editedFilename string)

	EditIssue(issue *api.Issue) (editedIssue *api.Issue)
}
