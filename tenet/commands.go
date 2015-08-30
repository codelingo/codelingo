package tenet

import (
	"encoding/json"
	"net/rpc/jsonrpc"
	"os"

	"github.com/juju/errors"
	"github.com/lingo-reviews/dev/tenet"
	"github.com/natefinch/pie"
)

type ReviewResult struct {
	TenetName string
	Issues    []*tenet.Issue
	Errs      []string
}

func (t *Tenet) Review(args ...string) (*ReviewResult, error) {
	// We need to create a new container for each Review (to mount /source
	// dir). TODO(waigani) We ignore any error. Though we should only be
	// ignoring a "not found" error.
	t.RemoveContainer()

	var result string
	err := t.call("Review", &result, args...)
	if err != nil {
		return nil, errors.Annotate(err, "error calling method Review")
	}

	reviewResult := &ReviewResult{}
	err = json.Unmarshal([]byte(result), reviewResult)
	reviewResult.TenetName = t.Name
	return reviewResult, err
}

func (t *Tenet) Help() (string, error) {
	var response string
	if err := t.call("Help", &response); err != nil {
		return "", err
	}
	return response, nil
}

func (t *Tenet) Version() (string, error) {
	var response string
	if err := t.call("Version", &response); err != nil {
		return "", err
	}
	return response, nil
}

func (t *Tenet) Debug(args ...string) string {
	var response string
	err := t.call("Debug", &response, args...)
	if err != nil {
		response += " error: " + err.Error()
	}
	return response
}

func (t *Tenet) CommentSet() (*tenet.CommentSet, error) {
	var comments tenet.CommentSet
	err := t.call("CommentSet", &comments)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &comments, nil
}

// TODO(matt) currently, call is hardcoded to call the tenet via a docker
// container. Make this configurable, with a "tenet-protocol" option, such
// that tenets can be executed as:
// - binary executables on host
// - remote web service
// - lxc
// - lxd
// - etc...
//
// result must be a pointer of type compatable with that returned by the remote method.
func (t *Tenet) call(method string, result interface{}, args ...string) error {
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
	defer client.Close()
	if err != nil {
		return errors.Annotate(err, "error running plugin")
	}
	return client.Call("Tenet."+method, args, result)
}
