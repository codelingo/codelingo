package driver

import (
	"net/rpc/jsonrpc"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/natefinch/pie"

	devTenet "github.com/lingo-reviews/dev/tenet"
)

// Binary is a tenet driver to execute binary tenets found in ~/.lingo/tenets/<repo>/<tenet>
type Binary struct {
	Common
}

func NewBinary(ctx *cli.Context, cfg Common) (*Binary, error) {
	return &Binary{
		Common: Common{
			Driver:  "binary",
			Name:    cfg.Name,
			Options: cfg.Options,
			context: ctx,
		},
	}, nil
}

// Check that a file exists at the expected location and is executable.
func (d *Binary) InitDriver() error {
	// TODO: implement
	return nil
}

// Do nothing - user is responsible for managing binary tenets.
func (d *Binary) Pull() error {
	return nil
}

// TODO: Use a makeCall proxy function to avoid repeating Help/Review/Version code across drivers
func (d *Binary) Review(args ...string) (*ReviewResult, error) {
	var result string
	err := d.call("Review", &result, args...)
	if err != nil {
		return nil, errors.Annotate(err, "error calling method Review")
	}

	return decodeResult(d.Name, result)
}

func (d *Binary) Help(args ...string) (string, error) {
	var response string
	if err := d.call("Help", &response, args...); err != nil {
		return "", err
	}
	return response, nil
}

func (d *Binary) Version() (string, error) {
	var response string
	if err := d.call("Version", &response); err != nil {
		return "", err
	}
	return response, nil
}

func (d *Binary) Debug(args ...string) string {
	var response string
	err := d.call("Debug", &response, args...)
	if err != nil {
		response += " error: " + err.Error()
	}
	return response
}

func (d *Binary) CommentSet() (*devTenet.CommentSet, error) {
	var comments devTenet.CommentSet
	err := d.call("CommentSet", &comments)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &comments, nil
}

// result must be a pointer of type compatable with that returned by the remote method.
func (d *Binary) call(method string, result interface{}, args ...string) error {
	// TODO: Should names like lingoHomeFlg be public so we can use them here/anywhere instead of hard-coding 'lingo-home'?
	tenetPath := path.Join(d.context.GlobalString("lingo-home"), "tenets", d.String())

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, tenetPath)
	if err != nil {
		return errors.Annotate(err, "error running tenet")
	}
	defer client.Close()

	return client.Call("Tenet."+method, args, result)
}
