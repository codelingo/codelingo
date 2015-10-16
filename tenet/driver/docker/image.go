package docker

import (
	"log"

	"github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"
)

func apiAuth() (docker.AuthConfiguration, error) {
	auths, err := docker.NewAuthConfigurationsFromDockerCfg()
	if err != nil {
		return docker.AuthConfiguration{}, errors.Errorf("error getting auth config: %s", err.Error())
	}
	if auth, ok := auths.Configs["https://index.docker.io/v1/"]; ok {
		return auth, nil
	}
	// otherwise return first auth found
	for _, auth := range auths.Configs {
		return auth, nil
	}
	return docker.AuthConfiguration{}, errors.New("auth not found")
}

func PullImage(client *docker.Client, name string, registry string, tag string) error {
	opts := docker.PullImageOptions{
		Repository: name,
		Registry:   registry,
		Tag:        tag,
	}
	auth, err := apiAuth()
	if err != nil {
		// just log err. We should be able to pull without auth.
		log.Printf("error getting auth config: %s", err.Error())
	}
	return client.PullImage(opts, auth)
}

// HaveImage returns true if the tenet image can be found locally.
func HaveImage(client *docker.Client, taggedName string) bool {
	img, _ := client.InspectImage(taggedName)
	return img != nil
}
