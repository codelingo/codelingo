package tenet

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/lingo-reviews/dev/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/lingo-reviews/lingo/tenet/driver"
	"github.com/lingo-reviews/lingo/tenet/service"
)

type Tenet interface {
	// Pull uses the tenet's driver to pull its image from its repository.
	Pull(bool) error

	// Service initiates the backing mirco-service and returns an object
	// to interact with it. Note the service needs to be explicitly
	// started and stopped.
	Service() (TenetService, error)
}

type TenetService interface {

	// Start the mirco-service and establish a connection to it.
	Start() error

	// Stop the micro-service. This closes the connection to the service. If
	// it is a locally running micro-service, the locally running process will
	// be killed.
	Stop() error

	// Review reads from filesc and writes to issuesc. It is the
	// responsibility of the caller to close filesc. The stream to the service
	// will stay open until filesc is closed.
	Review(filesc <-chan string, issuesc chan<- *api.Issue) error

	// Info returns all metadata about this tenet.
	Info() (*api.Info, error)

	// Returns the language the tenet applys to.
	Language() (string, error)
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
		d = &driver.DryRun{b}
	case "binary":
		d = &driver.Binary{b}
	default:
		return nil, fmt.Errorf("Unknown driver specified: %q", b.Driver)
	}
	return &tenet{Driver: d, options: b.ConfigOptions}, nil
}

func (t *tenet) Service() (TenetService, error) {
	s, err := t.Driver.Service()
	if err != nil {
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

	return &tenetService{
		Service:      s,
		cfg:          cfg,
		editFilename: t.Driver.EditFilename,
		editIssue:    t.Driver.EditIssue,
	}, nil
}

// tenetService implements TenetService.
type tenetService struct {
	service.Service
	client       api.TenetClient
	conn         *grpc.ClientConn
	cfg          *api.Config
	editFilename func(string) string
	editIssue    func(*api.Issue) *api.Issue
	info         *api.Info
}

func (t *tenetService) Start() error {
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

// Stop closes the connection, if local, stops the backing service.
func (t *tenetService) Stop() error {
	grpclog.Println("closing conn")
	err := t.conn.Close()
	if err1 := t.Service.Stop(); err1 != nil {
		err = err1
	}
	return err
}

func (t *tenetService) Language() (string, error) {
	i, err := t.Info()
	if err != nil {
		return "", errors.Trace(err)
	}

	return i.Language, nil
}

// Review will block and close the issues chan once the review is complete. It
// should be run in a gorountine.
func (t *tenetService) Review(filesc <-chan string, issues chan<- *api.Issue) error {
	stream, err := t.client.Review(context.Background())
	if err != nil {
		return err
	}
	waitc := make(chan struct{})
	go func() {
		for {
			issue, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				grpclog.Println("failed to receive an issue : %v", err)
			}
			grpclog.Printf("Got issue %v", issue)
			issues <- t.editIssue(issue)
		}
	}()

	for filename := range filesc {
		file := &api.File{Name: t.editFilename(filename)}
		if err := stream.Send(file); err != nil {
			grpclog.Println("failed to send a file %q: %v", filename, err)
		}
		grpclog.Printf("sent file %q", filename)
	}

	<-waitc
	close(issues)
	stream.CloseSend()
	return nil
}

func (t *tenetService) Info() (*api.Info, error) {
	if t.info == nil {
		i, err := t.client.GetInfo(context.Background(), &api.Nil{})
		if err != nil {
			return nil, errors.Trace(err)
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
	if t, err := New(ctx, b); err != nil {
		return t, nil
	}

	b.Driver = "docker"
	if t, err := New(ctx, b); err != nil {
		return t, nil
	}

	return nil, errors.Errorf("No driver available for %s", name)
}

func RenderedDescription(t Tenet) (string, error) {
	s, err := t.Service()
	if err != nil {
		return "", errors.Trace(err)
	}
	if err := s.Start(); err != nil {
		return "", err
	}
	defer s.Stop()
	info, err := s.Info()
	if err != nil {
		return "", errors.Trace(err)
	}

	tpl, err := template.New("desc template").Parse(info.Description)
	if err != nil {
		return "", err
	}

	opts := make(map[string]string)
	for _, opt := range info.Options {
		opts[opt.Name] = opt.Value
	}

	var rendered bytes.Buffer
	if err = tpl.Execute(&rendered, opts); err != nil {
		return "", err
	}

	return rendered.String(), nil
}
