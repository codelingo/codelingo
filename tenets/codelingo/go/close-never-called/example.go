package main

import (
	"errors"
	"os"
)

func main() {
	if err := do(); err != nil {
		panic(err)
	}
}

func do() error {
	file, err := os.Open("some.file")
	if err != nil {
		return err
	}
	file.Close()
	if err := other(); err != nil {
		return err // File is never closed if error happens here
	}
	if err := file.Close(); err != nil {
		return err
	}

	return nil
}

func other() error {
	return errors.New("an error")
}
