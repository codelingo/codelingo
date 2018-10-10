package docs

import (
	"io/ioutil"
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

	staged, err := gitCMD("-C", workingDir, "ls-tree", "-r", "--name-only", "--full-tree", "HEAD")
	if err != nil {
		return nil, errors.Trace(err)
	}

	unstaged, err := gitCMD("-C", workingDir, "ls-files", "--others", "--exclude-standard")
	if err != nil {
		return nil, errors.Trace(err)
	}

	files := strings.Split(staged, "\n")
	files = append(files, strings.Split(unstaged, "\n")...)

	lingoFilepaths := []string{}
	for _, filepath := range files {
		if isLingoFile(filepath) {
			lingoFilepaths = append(lingoFilepaths, filepath)
		}
	}

	return lingoFilepaths, nil
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
