package hello_test

import (
	"testing"

	"github.com/codelingo/codelingo/tenets/codelingo/jenkinsx/parallel-in-tests"
)

func TestHelloCorrect(t *testing.T) {
	t.Parallel()
	hello.Hello("correct")
}

func TestHelloIncorrect(t *testing.T) { // ISSUE
	hello.Hello("incorrect")
}
