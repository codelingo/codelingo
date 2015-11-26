package tenet

import (
	"fmt"
	"io"
	"time"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/lingo-reviews/dev/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/lingo-reviews/dev/tenet/log"
	"github.com/lingo-reviews/lingo/tenet/driver"
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
	Review(filesc <-chan string, issuesc chan<- *api.Issue) error

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

// NewTenet takes a base tenet and builds and returns a Tenet with a backing
// driver and service.
func New(ctx *cli.Context, b *driver.Base) (Tenet, error) {
	var d driver.Driver
	switch b.Driver {
	case "docker", "": // Default driver
		d = &driver.Docker{Base: b}
	case "dryrun":
		return &dryRun{}, nil
	case "binary":
		d = &driver.Binary{b}
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
	}

	if err := s.start(); err != nil {
		log.Print("got err opening service")
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
	err := t.conn.Close()
	log.Println("stopping service")
	if err1 := t.Service.Stop(); err1 != nil {
		err = err1
	}
	return err
}

// Review sets up two goroutines 1. to send files to the service from filesc,
// the other to recieve issues from the service on issuesc.
func (t *tenetService) Review(filesc <-chan string, issuesc chan<- *api.Issue) error {
	stream, err := t.client.Review(context.Background())
	if err != nil {
		return err
	}

	// first setup our issues chan to read from the service.
	go func(issuesc chan<- *api.Issue) {
		for {
			log.Println("waiting for issues")
			issue, err := stream.Recv()
			if err == io.EOF {
				log.Println("no more issues from tenet")
				// Close our local issues channel.
				close(issuesc)
				return
			}
			if err != nil {
				log.Fatalln("failed to receive an issue : %v", err)
			}
			if issue == nil {
				log.Fatalln("Issue is nil, this should never happen")
				return
			}
			log.Println("got an issue")
			issuesc <- t.editIssue(issue)
		}
	}(issuesc)

	// next, setup a goroutine to send our files to the service to review.
	go func(filesc <-chan string) {
		for {
			select {
			case filename, ok := <-filesc:
				if !ok && filename == "" {
					log.Println("client filesc closed. Closing send stream.")
					// Close the file send stream.
					stream.CloseSend()
					return
				}
				file := &api.File{Name: t.editFilename(filename)}
				if err := stream.Send(file); err != nil {
					log.Println("failed to send a file %q: %v", filename, err)
				}
				log.Printf("sent file %q", filename)
			case <-time.After(10 * time.Second):
				log.Fatal("timed out waiting for a filename")
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

type Options map[string]interface{}

// Configure configures the tenet with options set in .lingo or on the
// CLI. Any options not specified will fallback to their default value.
func (t *tenetService) configure() error {
	_, err := t.client.Configure(context.Background(), t.cfg)
	return err
}

// Any returns an initialised tenet using any driver that is locally available.
func Any(ctx *cli.Context, name string, options map[string]interface{}) (Tenet, error) {
	b := &driver.Base{
		Name:          name,
		ConfigOptions: options,
		Driver:        "binary",
	}

	// Try drivers in order of failure speed
	if t, err := New(ctx, b); err == nil {
		return t, nil
	}

	b.Driver = "docker"
	if t, err := New(ctx, b); err == nil {
		return t, nil
	}

	return nil, errors.Errorf("No driver available for %s", name)
}
