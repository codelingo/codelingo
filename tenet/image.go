package tenet

import (
	"log"

	"github.com/fsouza/go-dockerclient"
)

func (t *Tenet) PullImage() error {
	opts := docker.PullImageOptions{
		Repository: t.Name,
		Registry:   t.Registry,
		Tag:        t.Tag,
	}
	auth, err := apiAuth()
	if err != nil {
		// just log err. We should be able to pull without auth.
		log.Printf("error getting auth config: %s", err.Error())
	}
	return t.dockerClient.PullImage(opts, auth)
}

// HaveImage returns true if the tenet image can be found locally.
func (t *Tenet) HaveImage() bool {
	img, _ := t.dockerClient.InspectImage(t.String())
	return img != nil
}
