package main

import (
	"fmt"
)

func doStuff() string {
	return "Things"
}

func sayStuff() {
	fmt.Println("Things")
}


func main() {

	defer doStuff() // Issue
	defer sayStuff() // Non Issue

	defer func() { // Non Issue
		fmt.Println("Hello")
	}()

	defer func() string { // Issue
		return "Hello"
	}()
}
