package drivers

import (
	"net/rpc/jsonrpc"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/juju/errors"
	"github.com/natefinch/pie"

	"github.com/lingo-reviews/dev/api"
	devTenet "github.com/lingo-reviews/dev/tenet"
	"github.com/lingo-reviews/lingo/tenet"
)

type Docker struct {
	context *cli.Context
}

func newDocker(c *cli.Context) (Driver, error) {
	return &Docker{
		context: c,
	}, nil
}

func (d *Docker) Review(t *tenet.Tenet, args ...string) (*tenet.ReviewResult, error) {
	// We need to create a new container for each Review (to mount /source
	// dir). TODO(waigani) We ignore any error. Though we should only be
	// ignoring a "not found" error.
	t.RemoveContainer() // TODO: Move this to docker driver only

	var result string
	err := d.call("Review", t, &result, args...)
	if err != nil {
		return nil, errors.Annotate(err, "error calling method Review")
	}

	return decodeResult(t.Name, result)
}

func (d *Docker) Help(t *tenet.Tenet, args ...string) (string, error) {
	var response string
	if err := d.call("Help", t, &response, args...); err != nil {
		return "", err
	}
	return response, nil
}

func (d *Docker) Version(t *tenet.Tenet) (string, error) {
	var response string
	if err := d.call("Version", t, &response); err != nil {
		return "", err
	}
	return response, nil
}

func (d *Docker) Debug(t *tenet.Tenet, args ...string) string {
	var response string
	err := d.call("Debug", t, &response, args...)
	if err != nil {
		response += " error: " + err.Error()
	}
	return response
}

func (d *Docker) CommentSet(t *tenet.Tenet) (*devTenet.CommentSet, error) {
	var comments devTenet.CommentSet
	err := d.call("CommentSet", t, &comments)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &comments, nil
}

func (d *Docker) call(method string, t *tenet.Tenet, result interface{}, args ...string) error {
	containerName := t.ContainerName()

	// reuse existing container
	dockerArgs := []string{"start", "-i", containerName}

	if !t.HaveContainer() {
		// start new container
		dockerArgs = []string{"run", "-i"}
		if method == "Review" {
			// mount pwd as read only dir at root of container
			pwd, err := os.Getwd()
			if err != nil {
				return errors.Trace(err)
			}
			dockerArgs = append(dockerArgs, "-v", pwd+":/source:ro")
		}
		dockerArgs = append(dockerArgs, "--name", containerName, t.String())
	}

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, "docker", dockerArgs...)
	if err != nil {
		return errors.Annotate(err, "error running tenet")
	}
	defer client.Close()

	pass := len(args) // Number of args to leave alone - for non-review commands: all of them
	// Prepend filenames with /source/ for reviews
	if method == "Review" {
		cmd := api.ReviewCMD()
		err := cmd.Flags.Parse(args)
		if err != nil {
			return errors.Annotate(err, "could not parse arguments")
		}
		pass -= cmd.Flags.NArg()
	}

	var transformedArgs []string
	for i, a := range args {
		if i < pass {
			transformedArgs = append(transformedArgs, a)
		} else {
			transformedArgs = append(transformedArgs, path.Join("/source/", a))
		}
	}

	return client.Call("Tenet."+method, transformedArgs, result)
}
