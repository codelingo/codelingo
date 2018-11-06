package main

import "fmt"
import "errors"

func tryTheThing() (string, error) {
	return "hi", nil
}

var ErrDidNotWork = errors.New("did not work")

func DoTheThing(reallyDoIt bool) (err error) {
  if reallyDoIt {
    result, err := tryTheThing()
    if err != nil || result != "it worked" {
      err = ErrDidNotWork
    }
  }
  return err
}

func main() {
	e := DoTheThing(true)
	fmt.Printf("%v\n", e)
	// e.g. result, err := tryTheThing(). look for err := within a function where err is a parameter
}