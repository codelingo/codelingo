package util

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	goDocker "github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"
)

// TODO(anyone): Change this back to '.lingo' after making config loader check if
//               .lingo is file (not dir) before reading.
const (
	defaultHome = ".lingo_home"
)

func OpenFileCmd(editor, filename string, line int64) (*exec.Cmd, error) {
	app, err := exec.LookPath(editor)
	if err != nil {
		return nil, err
	}

	switch editor {
	case "subl", "sublime":
		return exec.Command(app, fmt.Sprintf("%s:%d", filename, line)), nil
		// TODO(waigani) other editors?
		// TODO(waigani) make the format a config var
	}

	// Making this default as vi, vim, nano, emacs all do it this way. These
	// are all terminal apps, so take over stdout etc.
	cmd := exec.Command(app, fmt.Sprintf("+%d", line), filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd, nil
}

func MustLingoHome() string {
	lHome, err := LingoHome()
	if err != nil {
		panic(err)
	}
	return lHome
}

func LingoHome() (string, error) {
	if lHome := os.Getenv("LINGO_HOME"); lHome != "" {
		return lHome, nil
	}
	home, err := UserHome()
	if err != nil {
		return "", err
	}

	return path.Join(home, defaultHome), nil
}

func UserHome() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func LingoBin() (string, error) {
	if bin := os.Getenv("LINGO_HOME"); bin != "" {
		return bin, nil
	}

	lHome, err := LingoHome()
	if err != nil {
		return "", errors.Trace(err)
	}

	return filepath.Join(lHome, "tenets"), nil

}

func BinTenets() ([]string, error) {
	binDir, err := LingoBin()
	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(binDir + "/*/*")
	if err != nil {
		return nil, err
	}

	tenets := make([]string, len(files))
	for i, f := range files {
		f = strings.TrimPrefix(f, binDir+"/")
		tenets[i] = f
	}
	return tenets, nil
}

func DockerTenets() ([]string, error) {
	// TODO(waigani) each image needs to be build with reviews.lingo.name for
	// this to work. We are using the label because it seems the client can't
	// find the image name.

	// c, err := DockerClient()
	// if err != nil {
	// 	return nil, errors.Trace(err)
	// }

	// images, err := c.ListImages(goDocker.ListImagesOptions{})
	// if err != nil {
	// 	return nil, errors.Trace(err)
	// }

	// tenets
	// for _, i := range images {
	// 	if l, ok := i.Labels["reviews.lingo.name"]; ok {
	// 		xxx.Print(l)
	// 	}
	// }

	return nil, nil
}

func DockerClient() (*goDocker.Client, error) {
	// TODO(waigani) get endpoint from ~/.lingo/config.toml
	endpoint := "unix:///var/run/docker.sock"
	return goDocker.NewClient(endpoint)
}
