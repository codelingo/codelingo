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

func writeSafeField(t thing) {
	t.name = "boat"
}

// Concurrent read and writes to a field of a copy of an instance of thing. No Issue
func concurrentCopyReadWrites() {

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

// Concurrent reads to a field of a copy of an instance of thing. No Issue
func concurrentCopyReads() {

	t := thing{"Object"}

	readSafeField(t)

	go readSafeField(t)

	go func(t thing) {
		fmt.Println(t.name)
	}(t)
}

// Concurrent writes to a field of a copy of an instance of thing. No Issue
func concurrentCopyWrites() {

	t := thing{"Object"}

	writeSafeField(t)

	go writeSafeField(t)

	go func(t thing) {
		t.name = "plane"
	}(t)
}

// Concurrent reads to a field of a pointer to an instance of thing. No Issue
func concurrentPointerReads() {

	t := &thing{"Pointer"}

	readVulnerableField(t)

	go readVulnerableField(t)

	go func(t *thing) {
		fmt.Println(t.name)
	}(t)
}

// Concurrent writes to a field of a pointer to an instance of thing. Issue
func concurrentPointerWrites() {

	t := &thing{"Pointer"}

	go writeVulnerableField(t)

	go func(t *thing) {
		t.name = "bus"
	}(t)
}

// Concurrent reads and writes to a field of a pointer to an instance of thing. Issue
func concurrentPointerReadWrites() {

	t := &thing{"Pointer"}

	go readVulnerableField(t)

	go writeVulnerableField(t)

	go func(t *thing) {
		fmt.Println(t.name)
	}(t)

	go func(t *thing) {
		t.name = "car"
	}(t)

}

func main() {
	concurrentCopyReads()
	concurrentCopyWrites()
	concurrentCopyReadWrites()
	concurrentPointerReads()
	concurrentPointerWrites()
	concurrentPointerReadWrites()
}
