package main

import (
	"fmt"
)

type thing struct {
	name string
}

func readFromPointer(t *thing) {
	fmt.Println(t.name)
}

func safeRead(t thing) {
	fmt.Println(t.name)
}

func writeToPointer(t *thing) {
	t.name = "car"
}

// safe demonstrates what the tenet has no interest in, here we perform several
// reads and writes, some of them on different threads, but all directly to an
// instance of a 'thing' so there are no issues.
func safe() {

	t := thing{"Object"}

	safeRead(t)

	go safeRead(t)

	go func(t thing) {
		fmt.Println(t.name)
	}(t)

	go func(t thing) {
		t.name = "bus"
	}(t)
}

// unsafe also performs several reads and writes, but here we are using a pointer
// to a 'thing'. When this happens on a seperate thread to the calling context we
// have the potential for unexpected bahviour.
func unsafe() {

	t := &thing{"Pointer"}

	readFromPointer(t)

	go readFromPointer(t)

	go writeToPointer(t) // Issue

	go func(t *thing) {
		fmt.Println(t.name)
	}(t)

	go func(t *thing) { // Issue
		t.name = "car"
	}(t)

}

func main() {
	safe()
	unsafe()
}
