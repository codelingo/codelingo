package tenet

import (
	"log"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"
)

// TODO(matt) These will turn into very long container names. Is there a limit
// on docker container names? We should create shorter unique names.
func (t *Tenet) ContainerName() string {
	r := strings.NewReplacer("/", "_", ":", "_")
	return r.Replace(t.String()) + "_lingotenet"
}

func (t *Tenet) CreateContainer() (*docker.Container, error) {
	opts := docker.CreateContainerOptions{
		Name: t.ContainerName(),
	}
	return t.dockerClient.CreateContainer(opts)
}

func (t *Tenet) RemoveContainer() error {
	opts := docker.RemoveContainerOptions{
		ID: t.ContainerName(),
	}
	return t.dockerClient.RemoveContainer(opts)
}

func (t *Tenet) StopContainer() error {
	return t.dockerClient.StopContainer(t.ContainerName(), 5)
}

func (t *Tenet) StartContainer() error {
	id, err := t.ContainerID()
	if err != nil {
		return errors.Trace(err)
	}
	return t.dockerClient.StartContainer(id, &docker.HostConfig{})
}

func (t *Tenet) ContainerID() (string, error) {
	name := t.ContainerName()
	// We track the container by name instead of ID, as ID will change if
	// container is destroyed and recreated.
	opts := docker.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"name": {name},
		},
	}
	containers, err := t.dockerClient.ListContainers(opts)
	if err != nil {
		return "", err
	}
	switch {
	case len(containers) == 1:
		return containers[0].ID, nil
	case len(containers) > 1:
		return "", errors.Errorf("found more than one container for %q. This should never happen.", name)
	}
	return "", errors.Errorf("container %q not found", name)
}

func (t *Tenet) HaveContainer() bool {
	id, err := t.ContainerID()
	if err != nil {
		log.Printf(err.Error())
	}
	return id != ""
}
