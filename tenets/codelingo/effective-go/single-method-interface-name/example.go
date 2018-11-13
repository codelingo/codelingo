// Package main for tenet testing
package main

import (
	"fmt"
)

// Reader is an example interface
type Reader interface {
	Read() []byte
}

// SomeName is an example interface
type SomeName interface {
	Read() []byte
}

// Writer is an example interface
type Writer interface {
	Write() int
}

// SomeWriter is an example interface
type SomeWriter interface {
	Write() int
}

// DoesntCount is an example interface
type DoesntCount interface {
	Write() int
	Read() []byte
}

// Read is an example interface
type Read interface {
	Reader() []byte
}

func main() {
	fmt.Println("Hello, world!")
}
