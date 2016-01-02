package driver

import (
	"fmt"
	"path"
	"strings"

	goDocker "github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"

	"github.com/lingo-reviews/lingo/tenet/driver/docker"
	"github.com/lingo-reviews/lingo/util"
	"github.com/lingo-reviews/tenets/go/dev/api"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
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

	haveImage := docker.HaveImage(dClient, d.Name)

	fmt.Printf("\npulling %s ... ", d.Name)
	if haveImage && !update {
		fmt.Printf("%s has already been pulled. Use --update to update.\n", d.Name)
		return nil
	}

	if update || !haveImage {
		if err := docker.PullImage(dClient, d.Name, d.Registry, d.Tag); err != nil {
			return err
		}
	}
	fmt.Printf("done.\n")
	return nil
}

func (d *Docker) getDockerClient() (*goDocker.Client, error) {
	if d.dockerClient == nil {
		dClient, err := util.DockerClient()
		if err != nil {
			return nil, errors.Trace(err)
		}
		d.dockerClient = dClient
	}

	return d.dockerClient, nil
}

// Init the service.
func (d *Docker) Service() (Service, error) {
	log.Print("Docker Service called")

	// Ensure we have an image.
	c, err := d.getDockerClient()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if !docker.HaveImage(c, d.Name) {

		// TODO(waigani) I think we should ask for user confirmation.
		fmt.Printf("\nno local image found for %s. Pulling new image from %s", d.Name, d.Registry)
		if err := d.Pull(false); err != nil {
			return nil, err
		}
	}

	return docker.NewService(d.Name)
}

// Docker mounts source code under /source/ so we need to prepend this to all
// file names.
func (d *Docker) EditFilename(f string) string {
	return path.Join("/source/", f)
}

func (d *Docker) EditIssue(issue *api.Issue) (editedIssue *api.Issue) {
	start := issue.Position.Start.Filename
	end := issue.Position.End.Filename

	issue.Position.Start.Filename = strings.TrimPrefix(start, "/source/")
	issue.Position.End.Filename = strings.TrimPrefix(end, "/source/")

	return issue
}
