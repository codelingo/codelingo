package commands

import (
	"path/filepath"

	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/util"

	"os"
	"os/exec"
)

type binaryBuildCfg struct {
	*baseBuildCfg
}

// BuildGo builds and installs a binary tenet. It assumes go is installed and
// the we are in the root of the tenet package.
func (cfg *binaryBuildCfg) BuildGo() error {
	cfg.dw.start.Done()
	defer cfg.dw.bar.Increment()

	bin, err := util.LingoBin()
	if err != nil {
		return errors.Trace(err)
	}
	tenetPath := filepath.Join(bin, cfg.Owner, cfg.Name)
	util.Printf("Building Go binary: %s\n", tenetPath)

	cmd := exec.Command("go", "get", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cfg.dir
	if err := cmd.Run(); err != nil {
		return errors.Trace(err)
	}

	cmd2 := exec.Command("go", "build", "-o", tenetPath)
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	cmd2.Dir = cfg.dir

	// TODO(waigani) capture this return as build error at end.
	cmd2.Stderr = os.Stderr
	return cmd2.Run()
}
