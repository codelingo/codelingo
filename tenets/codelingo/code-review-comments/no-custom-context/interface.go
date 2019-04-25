package main

import "time"

// Context is a custom context interface
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
	String() string
}

func main() {
}
