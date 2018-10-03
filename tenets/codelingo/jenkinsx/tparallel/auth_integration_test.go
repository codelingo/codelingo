package tparallel_test

import (
	"testing"

	"github.com/codelingo/codelingo/tenets/codelingo/jenkinsx/tparallel"
)

func TestAuthConfigParallel(t *testing.T) {
	t.Parallel()
	auth.DoAuthStuff()
}

func TestAuthConfig(t *testing.T) {
	auth.DoAuthStuff()
}
