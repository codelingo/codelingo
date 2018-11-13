// Package main for tenet testing
package main

import (
	"fmt"
)

type Reader interface {
	Read() []byte
}

type SomeName interface {
	Read() []byte
}

type Writer interface {
	Write() int
}

type SomeWriter interface {
	Write() int
}

type DoesntCount interface {
	Write() int
	Read() []byte
}

type Read interface {
	Reader() []byte
}

func main() {
	fmt.Println("Hello, world!")
}
