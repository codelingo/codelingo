package main

import (
	"errors"
	"fmt"
)

func main() {
	s := 3
	if err := checked(s); err != nil {
		fmt.Println(err.Error())
	}
	unchecked(s)
}

func unchecked(in interface{}) {
	s := in.(string) // ISSUE
	fmt.Println(s)
}

func checked(in interface{}) error {
	s, ok := in.(string) // No issue
	if !ok {
		return errors.New("failed to assert type")
	}

	fmt.Println(s)
	return nil
}
