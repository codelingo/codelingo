package main

import (
	"fmt"
)

type thing struct {
	name string
}

func readVulnerableField(t *thing) {
	fmt.Println(t.name)
}

func readSafeField(t thing) {
	fmt.Println(t.name)
}

func writeVulnerableField(t *thing) {
	t.name = "car"
}

func (t *thing) readFromPointer() {
	fmt.Println(t.name)
}

func (t *thing) writeToPointer(name string) {
	t.name = name
}

// safe demonstrates what the tenet has no interest in, here we perform several
// reads and writes, some of them on different threads, but all directly to an
// instance of a 'thing' so there are no issues.
func safe() {

	t := thing{"Object"}

	readSafeField(t)

	go readSafeField(t)

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

	readVulnerableField(t)

	go readVulnerableField(t)

	go writeVulnerableField(t) // Issue

	go func(t *thing) {
		fmt.Println(t.name)
	}(t)

	go func(t *thing) { // Issue
		t.name = "car"
	}(t)

	t.readFromPointer()

	t.writeToPointer("plane")

	go t.writeToPointer("boat") // Issue
}

func main() {
	safe()
	unsafe()
}
