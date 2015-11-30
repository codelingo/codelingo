package commands

import (
	"path/filepath"

	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/util"

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
	cmd := exec.Command("go", "build", "-o", tenetPath)
	cmd.Dir = cfg.dir
	return cmd.Run()
}
