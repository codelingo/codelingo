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
	Issues []*tenet.Issue
	Errs   []string
}

func (t *Tenet) Review(args ...string) (*ReviewResult, error) {
	// We need to create a new container for each Review (to mount /source
	// dir). TODO(waigani) We ignore any error. Though we should only be
	// ignoring a "not found" error.
	t.RemoveContainer()

	result, err := t.call("Review", args...)
	if err != nil {
		return nil, errors.Trace(err)
	}

	reviewResult := &ReviewResult{}
	err = json.Unmarshal([]byte(result), reviewResult)
	return reviewResult, err
}

func (t *Tenet) Help() (string, error) {
	return t.call("Help")
}

func (t *Tenet) Version() (string, error) {
	return t.call("Version")
}

func (t *Tenet) Debug(args ...string) string {
	str, err := t.call("Debug", args...)
	if err != nil {
		str += " error: " + err.Error()
	}
	return str
}

func (t *Tenet) call(method string, args ...string) (string, error) {
	// this block is just for debugging
	// _ = dockerArgs
	// path := "/home/jesse/go/src/github.com/lingo-reviews/dev/tenetseed/tenetseed"
	// client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, path)
	// if err != nil {
	// 	log.Fatalf("Error running plugin: %s", err)
	// }
	// defer client.Close()

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
				return "", errors.Trace(err)
			}
			dockerArgs = append(dockerArgs, "-v", pwd+":/source:ro")
		}
		dockerArgs = append(dockerArgs, "--name", containerName, t.String())
	}

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, "docker", dockerArgs...)
	defer client.Close()
	if err != nil {
		return "", errors.Annotate(err, "error running plugin")
	}

	var result string
	err = client.Call("Tenet."+method, args, &result)
	if err != nil {
		return "", errors.Annotatef(err, "error calling method %q", method)
	}
	return result, nil
}
