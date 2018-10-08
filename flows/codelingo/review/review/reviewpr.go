// The review package contains helper methods to create review requests to be sent to the bot
// endpoint layer, especially CLAIR.
package review

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/juju/errors"
)

type PROpts struct {
	Host     string
	HostName string
	Owner    string
	Name     string
	PRID     int
}

// ParsePR produces a set of PROpts to be sent to CLAIR in the bot endpoint layer
func ParsePR(urlStr string) (*PROpts, error) {
	// TODO: Try to parse urlStr as a request for each VCS host

	result, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var parts []string
	parts = strings.Split(strings.Trim(result.Path, "/"), "/")
	switch result.Host {
	case "github.com":
		opts, err := parseGithubPR(parts)
		return opts, errors.Trace(err)
	case "gitlab.com":
		opts, err := parseGitlabPR(parts)
		return opts, errors.Trace(err)
	case "bitbucket.com":
		opts, err := parseBitBucketPR(parts)
		return opts, errors.Trace(err)
	default:
		return nil, errors.Errorf("Unrecognised host %s", result.Host)
	}
}

func parseGithubPR(urlPath []string) (*PROpts, error) {
	if l := len(urlPath); l != 4 {
		return nil, errors.Errorf("Github pull request URL needs to be in the following format: https://github.com/<username>/<repo_name>/pull/<pull_number>")
	}

	n, err := strconv.Atoi(urlPath[3])
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &PROpts{
		Host:     "github.com",
		HostName: "github_com",
		Owner:    urlPath[0],
		Name:     urlPath[1],
		PRID:     n,
	}, nil
}

// TODO: implement
func parseBitBucketPR(urlStr []string) (*PROpts, error) {
	return nil, errors.New("BitBucket not supported.")
}

// TODO: implement
func parseGitlabPR(urlStr []string) (*PROpts, error) {
	return nil, errors.New("Gitlab not supported.")
}
