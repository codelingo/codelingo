package tenet

import (
	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"
)

type Tenet struct {
	// Name of the image
	Name string `toml:"name"`

	// Tag of the image
	Tag string `toml:"tag"`

	// Registry server to pull the image from
	Registry string `toml:"registry"`

	// Config options for tenet
	Options map[string]interface{} `toml:"options"`

	dockerClient *docker.Client
}

var dClient *docker.Client

func init() {
	// TODO(waigani) get endpoint from ~/.lingo/config.toml
	endpoint := "unix:///var/run/docker.sock"
	var err error
	dClient, err = docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
}

// NewTenet builds and returns a Tenet. name is a docker image name of form <repo>/<image>:[tag]
func New(name string) (*Tenet, error) {
	parts := strings.Split(name, ":")
	t := Tenet{
		Name:     parts[0],
		Registry: "hub.docker.com",
	}

	l := len(parts)
	switch {
	case l > 2:
		return nil, errors.Errorf("%q is wrong format")
	case l == 2:
		t.Tag = parts[1]
	}

	if err := t.DockerInit(); err != nil {
		return nil, errors.Trace(err)
	}
	return &t, nil
}

// DockerInit prepares the object to talk to the docker image backing it. It
// also pulls the image if missing.
func (t *Tenet) DockerInit() error {
	t.dockerClient = dClient

	if !t.HaveImage() {
		if err := t.PullImage(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tenet) String() string {
	if t.Tag != "" {
		return t.Name + ":" + t.Tag
	}
	return t.Name
}

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
