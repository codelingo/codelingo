package service

import (
	"strings"

	"github.com/codegangsta/cli"
)

func serviceEndpoint(serviceName string) string {
	baseServiceURL := "https://do.lingo.reviews/service"
	// baseServiceURL := "http://localhost:8080/service"
	return strings.Trim(baseServiceURL, "/") + "/" + strings.Trim(serviceName, "/")
}

var Services = []cli.Command{
	reviewboard,
}
