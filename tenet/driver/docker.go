package driver

import (
	"os"
	"path"

	goDocker "github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"

	"github.com/lingo-reviews/lingo/tenet/driver/docker"
	"github.com/lingo-reviews/lingo/tenet/service"
)

// Docker is a tenet driver which runs tenets inside a docker container.
type Docker struct {
	*Base

	dockerClient *goDocker.Client
}

// Pull the image for this tenet from the given registry.
func (d *Docker) Pull(update bool) error {
	dClient, err := d.getDockerClient()
	if err != nil {
		return errors.Trace(err)
	}

	if update || !docker.HaveImage(dClient, d.Name) {
		return docker.PullImage(dClient, d.Name, d.Registry, d.Tag)
	}
	return nil
}

func (d *Docker) getDockerClient() (*goDocker.Client, error) {
	if d.dockerClient == nil {
		// TODO(waigani) get endpoint from ~/.lingo/config.toml
		endpoint := "unix:///var/run/docker.sock"

		dClient, err := goDocker.NewClient(endpoint)
		if err != nil {
			return nil, err
		}
		d.dockerClient = dClient
	}

	return d.dockerClient, nil
}

// Init the service. Note: the service needs to be started and stopped.
func (d *Docker) Service() (service.Service, error) {
	dClient, err := d.getDockerClient()
	if err != nil {
		return nil, errors.Trace(err)
	}
	// We need to create a new container for each Review (to mount /source
	// dir). TODO(waigani) We ignore any error. Though we should only be
	// ignoring a "not found" error.
	docker.RemoveContainer(dClient, d.Name)

	// TODO(waigani)  TECHDEBT we always start a new container. Once tenets gain state, reuse same container.

	containerName := docker.ContainerName(d.Name)
	// TODO(waigani) TECHDEBT this is all we need once we reuse existing container.
	// dockerArgs := []string{"start", "-i", containerName}

	// start new container
	pwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Trace(err)
	}
	dockerArgs := []string{
		// start a new container
		"run", "-i",

		// mount pwd as read only dir at root of container
		"-v", pwd + ":/source:ro",

		"--name", containerName, d.Name,
	}

	return service.NewLocal("docker", dockerArgs...), nil
}

// Docker mounts source code under /source/ so we need to prepend this to all
// file names.
func (d *Docker) EditFilename(f string) string {
	return path.Join("/source/", f)
}
