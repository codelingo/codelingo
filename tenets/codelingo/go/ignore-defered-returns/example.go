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
}
