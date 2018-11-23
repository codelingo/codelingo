package main

import (
	"fmt"

	"github.com/juju/errors"
)

func main() {
	_ = doStuffA("some param")
	_ = doStuffB("some param")
}

func doStuffA(param string) error {
	// bad
	return errors.New(fmt.Sprintf("Don't call with \"%s\" param - it will literally do nothing!", param))
}

func doStuffN(param string) error {
	// good
	return errors.Errorf("Don't call with \"%s\" param - it will literally do nothing!", param)
}
