package main

import (
	"fmt"

	"github.com/juju/errors"
)

func main() {
	_ = doStuff("some param")
}

func doStuff(param string) error {
	return errors.New(fmt.Sprintf("Don't call with \"%s\" param - it will literally do nothing!", param))
}

func paramErr(param string) error {
	return errors.New(fmt.Sprintf("parameter error: \"%s\"", param))
}

func doOtherStuff(param string) error {
	return errors.Errorf("Don't call with \"%s\" param - it will literally do nothing!", param)
}
