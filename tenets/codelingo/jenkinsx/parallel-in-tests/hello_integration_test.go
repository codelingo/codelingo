package hello_test

import (
	"testing"

	"github.com/codelingo/codelingo/tenets/codelingo/jenkinsx/parallel-in-tests"
)

func TestHelloIntegration(t *testing.T) { // Shouldn't match
	hello.Hello("incorrect")
}
