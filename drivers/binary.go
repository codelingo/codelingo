package drivers

import (
	"net/rpc/jsonrpc"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/natefinch/pie"

	devTenet "github.com/lingo-reviews/dev/tenet"
	"github.com/lingo-reviews/lingo/tenet"
)

// Binary is a tenet driver to execute binary tenets found in ~/.lingo/tenets/<repo>/<tenet>
type Binary struct {
	context *cli.Context
}

func newBinary(c *cli.Context) (Driver, error) {
	return &Binary{
		context: c,
	}, nil
}

// TODO: What's the best way to share common code between drivers (polymorphic 'call' would be ideal)?
func (d *Binary) Review(t *tenet.Tenet, args ...string) (*tenet.ReviewResult, error) {
	var result string
	err := d.call("Review", t, &result, args...)
	if err != nil {
		return nil, errors.Annotate(err, "error calling method Review")
	}

	return decodeResult(t.Name, result)
}

func (d *Binary) Help(t *tenet.Tenet, args ...string) (string, error) {
	var response string
	if err := d.call("Help", t, &response, args...); err != nil {
		return "", err
	}
	return response, nil
}

func (d *Binary) Version(t *tenet.Tenet) (string, error) {
	var response string
	if err := d.call("Version", t, &response); err != nil {
		return "", err
	}
	return response, nil
}

func (d *Binary) Debug(t *tenet.Tenet, args ...string) string {
	var response string
	err := d.call("Debug", t, &response, args...)
	if err != nil {
		response += " error: " + err.Error()
	}
	return response
}

func (d *Binary) CommentSet(t *tenet.Tenet) (*devTenet.CommentSet, error) {
	var comments devTenet.CommentSet
	err := d.call("CommentSet", t, &comments)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &comments, nil
}

// result must be a pointer of type compatable with that returned by the remote method.
func (d *Binary) call(method string, t *tenet.Tenet, result interface{}, args ...string) error {
	// TODO: Should names like lingoHomeFlg be public so we can use them here/anywhere?
	tenetPath := path.Join(d.context.GlobalString("lingo-home"), "tenets", t.String())

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, tenetPath)
	if err != nil {
		return errors.Annotate(err, "error running tenet")
	}
	defer client.Close()

	return client.Call("Tenet."+method, args, result)
}
