package tenet

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/lingo-reviews/tenets/go/dev/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	tomb "gopkg.in/tomb.v1"

	"github.com/lingo-reviews/lingo/tenet/driver"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
)

type Tenet interface {
	// Pull uses the tenet's driver to pull its image from its repository.
	Pull(bool) error

	// Service initiates the backing mirco-service and returns an object
	// to interact with it. Note the service needs to be explicitly
	// started and stopped.
	OpenService() (TenetService, error)
}

type TenetService interface {

	// Close stops the micro-service. This closes the connection to the service. If
	// it is a locally running micro-service, the locally running process will
	// be killed.
	Close() error

	// Review reads from filesc and writes to issuesc. It is the
	// responsibility of the caller to close filesc. The stream to the service
	// will stay open until filesc is closed.
	Review(filesc <-chan *api.File, issuesc chan<- *api.Issue, t *tomb.Tomb) error

	// Info returns all metadata about this tenet.
	Info() (*api.Info, error)

	// start the backing service process
	start() error
}

// tenet implenets Tenet.
type tenet struct {
	driver.Driver
	options map[string]interface{}
}

// New takes a base tenet config and builds and returns a Tenet with a backing
// driver.
func New(ctx *cli.Context, b *driver.Base) (Tenet, error) {
	var d driver.Driver
	switch b.Driver {
	case "docker", "": // Default driver
		d = &driver.Docker{Base: b}
	case "dryrun":
		return &dryRun{}, nil
	case "binary":
		d = &driver.Binary{Base: b}
	default:
		return nil, fmt.Errorf("Unknown driver specified: %q", b.Driver)
	}
	return &tenet{Driver: d, options: b.ConfigOptions}, nil
}

func (t *tenet) OpenService() (TenetService, error) {
	ds, err := t.Driver.Service()
	if err != nil {
		log.Print(err.Error()) // TODO(waigani) this logs are quick hacks. Work out the error paths and log them all at the root.
		return nil, errors.Trace(err)
	}

	cfg := &api.Config{}
	for k, v := range t.options {
		cfg.Options = append(cfg.Options,
			&api.Option{
				Name:  k,
				Value: fmt.Sprintf("%v", v),
			})
	}

	s := &tenetService{
		Service:      ds,
		cfg:          cfg,
		editFilename: t.Driver.EditFilename,
		editIssue:    t.Driver.EditIssue,
		mutex:        &sync.Mutex{},
	}

	if err := s.start(); err != nil {
		log.Println("got err opening service")
		// TODO(waigani) add retry logic here. 1. Keep retrying until service
		// is up. 2. Keep retrying until service is connected.
		log.Printf("err: %#v", errors.ErrorStack(err))

		return nil, errors.Trace(err)
	}
	log.Print("opened service, no issue")
	return s, nil
}

// tenetService implements TenetService.
type tenetService struct {
	driver.Service
	client       api.TenetClient
	conn         *grpc.ClientConn
	cfg          *api.Config
	editFilename func(string) string
	editIssue    func(*api.Issue) *api.Issue
	info         *api.Info
	mutex        *sync.Mutex
}

func (t *tenetService) start() error {
	if t.Service == nil {
		return errors.New("service is nil. Has the tenet driver been initialized?")
	}
	s := t.Service
	if err := s.Start(); err != nil {
		return errors.Trace(err)
	}
	var err error
	t.conn, err = s.DialGRPC()
	if err != nil {
		return errors.Trace(err)
	}
	c := api.NewTenetClient(t.conn)
	t.client = c
	return t.configure()
}

// Close closes the connection and, if local, stops the backing service.
func (t *tenetService) Close() error {
	log.Println("closing conn")

	if t.Service == nil {
		return errors.New("attempted to close a nil service. Has it been started?")
	}
	err := t.conn.Close()
	log.Println("stopping service")
	if err1 := t.Service.Stop(); err1 != nil {
		err = err1
	}
	return err
}

// type closerHelper struct {
// 	closable interface{}
// }

// func newCloser(i interface{}) *closerHelper {
// 	return
// }

// func (c *closerHelper) Close() {
// 	if c.closable != nil {
// 		close(c.closable)
// 		c.closable = nil
// 	}
// }

// Review sets up two goroutines. One to send files to the service from filesc,
// the other to recieve issues from the service on issuesc.
func (t *tenetService) Review(filesc <-chan *api.File, issuesc chan<- *api.Issue, filesTM *tomb.Tomb) error {
	stream, err := t.client.Review(context.Background())
	if err != nil {
		return err
	}

	// first setup our issues chan to read from the service.
	go func(issuesc chan<- *api.Issue) {
		for {
			log.Println("waiting for issues")
			issue, err := stream.Recv()
			if err != nil {
				if err == io.EOF ||
					grpc.ErrorDesc(err) == "transport is closing" ||
					err.Error() == "timed out waiting for issues" { // TODO(waigani) error type
					log.Println("closing issuesc")
					// Close our local issues channel.
					close(issuesc)
					return
				}

				// TODO(waigani) in what error cases should we close issuesc?
				// Any other err we keep calm and carry on.
				log.Println("ERROR receiving an issue : %s", err.Error())
				continue
			}

			issuesc <- t.editIssue(issue)
		}
	}(issuesc)

	// next, setup a goroutine to send our files to the service to review.
	go func(filesc <-chan *api.File) {
		for {
			select {
			case file, ok := <-filesc:
				if !ok && file == nil {
					log.Println("client filesc closed. Closing send stream.")
					// Close the file send stream.
					err := stream.CloseSend()
					if err != nil {
						log.Println(err.Error())
					}
					return
				}

				file.Name = t.editFilename(file.Name)
				if err := stream.Send(file); err != nil {
					log.Println("failed to send a file %q: %v", file.Name, err)
				}
				log.Printf("sent file %q\n", file.Name)

				// Each tenet has a 5 second idle time. If we don't find any
				// files to send it in that time, we close this tenet down.
			case <-time.After(5 * time.Second):
				// this will close this instance of the tenet.
				filesTM.Kill(errors.New("timed out waiting for a filename"))

				return
			}
		}
	}(filesc)

	return nil
}

func (t *tenetService) Info() (*api.Info, error) {
	if t.info == nil {
		i, err := t.client.GetInfo(context.Background(), &api.Nil{})
		if err != nil {
			return nil, errors.Trace(err)
		}

		if i.Version == "" {
			i.Version = "0.0.0"
		}
		t.info = i
	}
	return t.info, nil
}

// Configure configures the tenet with options set in .lingo or on the
// CLI. Any options not specified will fallback to their default value.
func (t *tenetService) configure() error {
	_, err := t.client.Configure(context.Background(), t.cfg)
	return err
}
