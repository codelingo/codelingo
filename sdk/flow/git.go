package flow

import (
	"os/exec"

	"github.com/juju/errors"
)

func GitCMD(args ...string) (out string, err error) {
	cmd := exec.Command("git", args...)
	b, err := cmd.CombinedOutput()
	out = string(b)
	return out, errors.Annotate(err, out)
}
