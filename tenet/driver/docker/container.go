package docker

import (
	"log"
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"
)

// TODO: There's some dead code in this file (StartContainer, StopContainer), do we need it?

// TODO(matt) These will turn into very long container names. Is there a limit
// on docker container names? We should create shorter unique names.
func ContainerName(name string) string {
	r := strings.NewReplacer("/", "_", ":", "_")
	return r.Replace(name) + "_lingotenet"
}

func CreateContainer(client *docker.Client, name string) (*docker.Container, error) {
	opts := docker.CreateContainerOptions{
		Name: ContainerName(name),
	}
	return client.CreateContainer(opts)
}

func RemoveContainer(client *docker.Client, name string) error {
	opts := docker.RemoveContainerOptions{
		ID: ContainerName(name),
	}
	return client.RemoveContainer(opts)
}

func StopContainer(client *docker.Client, name string) error {
	return client.StopContainer(ContainerName(name), 5)
}

func StartContainer(client *docker.Client, name string) error {
	id, err := ContainerID(client, name)
	if err != nil {
		return errors.Trace(err)
	}
	return client.StartContainer(id, &docker.HostConfig{})
}

func ContainerID(client *docker.Client, name string) (string, error) {
	cName := ContainerName(name)
	// We track the container by name instead of ID, as ID will change if
	// container is destroyed and recreated.
	opts := docker.ListContainersOptions{
		All: true,
		Filters: map[string][]string{
			"name": {cName},
		},
	}
	containers, err := client.ListContainers(opts)
	if err != nil {
		return "", err
	}
	switch {
	case len(containers) == 1:
		return containers[0].ID, nil
	case len(containers) > 1:
		return "", errors.Errorf("found more than one container for %q. This should never happen.", cName)
	}
	return "", errors.Errorf("container %q not found", name)
}

func HaveContainer(client *docker.Client, name string) bool {
	id, err := ContainerID(client, name)
	if err != nil {
		log.Printf(err.Error())
	}
	return id != ""
}
