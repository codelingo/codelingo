package driver

import (
	"github.com/lingo-reviews/dev/api"
	"google.golang.org/grpc"
)

type Driver interface {
	Pull(bool) error

	Service() (Service, error)

	EditFilename(filename string) (editedFilename string)

	EditIssue(issue *api.Issue) (editedIssue *api.Issue)
}

// service handles operations with the underlying backing micro-service tenet
// server.
type Service interface {
	Start() error
	Stop() error
	// IsRunning() bool
	DialGRPC() (*grpc.ClientConn, error)
}
