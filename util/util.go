package util

import (
	"fmt"
	"os"
	"os/exec"
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
