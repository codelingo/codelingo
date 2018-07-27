package main

import "log"

// todo: all lower case

// tOdO: mixed case

// TODO: upper case

// todo no colon

// gasdTodo--

// main is the entry point to this example
func main() {
	err := serve()
	if err != nil {
		log.Fatal(err.Error()) // TODO: check error type before exiting
	}
}

// What about this comment?

// serve runs the server, returning an error if it crashes
func serve() error {
	// TODO: implement server
	return nil
}
