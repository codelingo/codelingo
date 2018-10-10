package docs

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
)

// Get all codelingo.yaml files that apply to the current review.
func GetLingoFiles(workingDir string) (map[string][]byte, error) {
	paths, err := getLingoFilepaths(workingDir)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// Only those paths that are either a child or parent of the working directory generate queries.
	// Go to the root of the repository to do a full test.
	lingoFiles := map[string][]byte{}
	for _, path := range paths {
		pathDir := filepath.Dir(path)
		if pathDir == "." {
			pathDir = ""
		}

		if strings.HasPrefix(pathDir, workingDir) || strings.HasPrefix(workingDir, pathDir) {

			fileSRC, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, errors.Trace(err)
			}

			lingoFiles[path] = fileSRC

		}
	}
	return lingoFiles, nil
}

// gets the file paths of all the codelingo.yaml files in the repo
func getLingoFilepaths(workingDir string) ([]string, error) {
	var files []string
	err := filepath.Walk(workingDir, func(path string, info os.FileInfo, err error) error {
		if isLingoFile(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Trace(err)
	}

	return files, nil
}

// isDotlingoFile returns if that given filepath has a recognised lingo extension.
func isLingoFile(file string) bool {

	filename := filepath.Base(file)
	return map[string]bool{
		"codelingo.yaml": true,
		"codelingo.yml":  true,
	}[filename]
}

// by any git command-line actions
func gitCMD(args ...string) (out string, err error) {
	cmd := exec.Command("git", args...)
	b, err := cmd.CombinedOutput()
	out = string(b)
	return out, errors.Annotate(err, out)
}
