package driver

import (
	"net/rpc/jsonrpc"
	"os"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	goDocker "github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"
	"github.com/natefinch/pie"

	"github.com/lingo-reviews/dev/api"
	"github.com/lingo-reviews/lingo/tenet/driver/docker"
)

// Docker is a tenet driver which runs tenets inside a docker container.
type Docker struct {
	Common

	DockerClient *goDocker.Client
}

func NewDocker(ctx *cli.Context, cfg Common) (*Docker, error) {
	parts := strings.Split(cfg.Name, ":")

	name := parts[0]
	tag := ""
	registry := cfg.Registry

	l := len(parts)
	switch {
	case l > 2:
		return nil, errors.Errorf("Tenet name '%q' is the wrong format")
	case l == 2:
		tag = parts[1]
	}

	if registry == "" {
		registry = "hub.docker.com"
	}

	return &Docker{
		Common: Common{
			Driver:   "docker",
			Name:     name,
			Tag:      tag,
			Registry: registry,
			Options:  cfg.Options,
			context:  ctx,
		},
	}, nil
}

// InitDriver prepares the tenet to talk to the docker image backing it. It
// also pulls the image if missing.
func (d *Docker) InitDriver() error {
	// TODO(waigani) get endpoint from ~/.lingo/config.toml
	endpoint := "unix:///var/run/docker.sock"

	dClient, err := goDocker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	d.DockerClient = dClient

	if err := d.Pull(); err != nil {
		return err
	}
	return nil
}

// Pull the image for this tenet from the given registry.
func (d *Docker) Pull() error {
	if !docker.HaveImage(d.DockerClient, d.Name) {
		return docker.PullImage(d.DockerClient, d.Name, d.Registry, d.Tag)
	}
	return nil
}

func (d *Docker) Review(args ...string) (*ReviewResult, error) {
	// We need to create a new container for each Review (to mount /source
	// dir). TODO(waigani) We ignore any error. Though we should only be
	// ignoring a "not found" error.
	docker.RemoveContainer(d.DockerClient, d.Name)

	var result string
	err := d.call("Review", &result, args...)
	if err != nil {
		return nil, errors.Annotate(err, "error calling method Review")
	}

	return decodeResult(d.Name, result)
}

func (d *Docker) Help(args ...string) (string, error) {
	var response string
	if err := d.call("Help", &response, args...); err != nil {
		return "", err
	}
	return response, nil
}

func (d *Docker) Description() (string, error) {
	var response string
	if err := d.call("Description", &response); err != nil {
		return "", err
	}
	return response, nil
}

func (d *Docker) Version() (string, error) {
	var response string
	if err := d.call("Version", &response); err != nil {
		return "", err
	}
	return response, nil
}

func (d *Docker) Debug(args ...string) string {
	var response string
	err := d.call("Debug", &response, args...)
	if err != nil {
		response += " error: " + err.Error()
	}
	return response
}

func (d *Docker) call(method string, result interface{}, args ...string) error {
	containerName := docker.ContainerName(d.Name)

	// reuse existing container
	dockerArgs := []string{"start", "-i", containerName}

	if !docker.HaveContainer(d.DockerClient, d.Name) {
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
		dockerArgs = append(dockerArgs, "--name", containerName, d.String())
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
