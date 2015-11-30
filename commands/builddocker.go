package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	goDocker "github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/util"
)

type dockerBuildCfg struct {
	OverwriteDockerfile bool `toml:"overwrite_dockerfile"`
	*baseBuildCfg
	dockerClient *goDocker.Client
}

// BuildGo builds and installs a binary tenet. It assumes go is installed.
func (cfg *dockerBuildCfg) BuildGo() error {
	cfg.dw.start.Done()
	defer cfg.dw.bar.Increment()

	dockerfile := filepath.Join(cfg.dir, "Dockerfile")
	var exists bool
	// Write a Dockerfile if one does not exist.
	if n, err := os.Stat(dockerfile); err != nil && !os.IsNotExist(err) {
		return errors.Trace(err)
	} else if n != nil {
		exists = true
	}

	if cfg.OverwriteDockerfile || !exists {
		err := cfg.writeDockerFile(dockerfile)
		if err != nil {
			return errors.Trace(err)
		}
	}

	c, err := cfg.client()
	if err != nil {
		return errors.Trace(err)
	}

	// TODO(waigani) remove when dev is published
	_, relDockerPath, lrRootDir, err := cfg.hackedPaths()
	if err != nil {
		return errors.Trace(err)
	}
	imageName := cfg.Owner + "/" + cfg.Name
	util.Printf("Building Docker image: %s\n", imageName)
	opts := goDocker.BuildImageOptions{
		Name:         imageName,
		Dockerfile:   relDockerPath, //dockerfile,
		ContextDir:   lrRootDir,     // "~/go/src/github.com/lingo-reviews/tenets/simpleseed", //  "~/go/src/github.com/lingo-reviews/tenets/simpleseed/", // TODO(waigani) when dev is published ContextDir =  cfg.dir
		OutputStream: os.Stdout,
	}

	return c.BuildImage(opts)
}

func (cfg *dockerBuildCfg) Publish() error {
	// c, err := cfg.client()
	// if err != nil {
	// 	return errors.Trace(err)
	// }

	// c.PushImage(opts, auth)

	return errors.New("not implemented")
}

// TODO(waigani) remove when dev is published
func lrRoot() (string, error) {
	h, err := util.UserHome()
	if err != nil {
		return "", errors.Trace(err)
	}

	return filepath.Join(h, "go", "src", "github.com", "lingo-reviews"), nil
}

// TODO(waigani) remove when dev is published
func (cfg *dockerBuildCfg) hackedPaths() (tenetRelPath, relDockerPath, lrRootDir string, err error) {
	// TODO(waigani) all this shit is because dev is not published
	lrRootDir, err = lrRoot()
	if err != nil {
		return "", "", "", errors.Trace(err)
	}
	tenetRelPath = strings.TrimPrefix(cfg.dir, lrRootDir+"/")
	relDockerPath = filepath.Join(tenetRelPath, "Dockerfile")
	return
}

// TODO(waigani) remove when dev is published
type dockerfileData struct {
	*baseBuildCfg
	TenetRoot string
}

func (cfg *dockerBuildCfg) writeDockerFile(dockerfilePath string) error {
	// TODO(waigani) remove when dev is published
	tenetRelPath, _, _, err := cfg.hackedPaths()
	if err != nil {
		return errors.Trace(err)
	}

	data := dockerfileData{
		baseBuildCfg: cfg.baseBuildCfg,
		TenetRoot:    tenetRelPath, // TODO(waigani) use cfg.dir when dev is published
	}
	out, err := util.FormatOutput(data, dockerfileTmpl)
	if err != nil {
		return errors.Trace(err)
	}

	dockerfilePath = filepath.Join(cfg.dir, "Dockerfile")
	return ioutil.WriteFile(dockerfilePath, []byte(out), 0644)
}

func (cfg *dockerBuildCfg) client() (*goDocker.Client, error) {
	if cfg.dockerClient == nil {
		var err error
		cfg.dockerClient, err = util.DockerClient()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return cfg.dockerClient, nil
}

var dockerfileTmpl = `
FROM golang
# FROM golang:onbuild # <---- This line is all that will be needed.

ENV LINGO_CONTAINER true

## ---------
# The following is only needed while lingo libs are privately hosted on
# bitbucket. Once they are published, 'FROM golang:onbuild' is all we need
# here. But for now we need to manually checkout the repos into the paths
# copied below.

COPY . /go/src/github.com/lingo-reviews
COPY {{.TenetRoot}} /go/src/app
WORKDIR /go/src/app
RUN go get -v -d
RUN go install -v
ENTRYPOINT /go/bin/app
## ----------

# This info is used for searching for tenet images.
LABEL reviews.lingo.name="{{.Owner}}/{{.Name}}" \
`[1:]
