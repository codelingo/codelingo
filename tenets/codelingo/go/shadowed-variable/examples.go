package main

import (
	"fmt"
)

func main() {
	foo()
	bar()
}

func foo() {
	x := 1
	fmt.Println(x)
	{
		x = 5  // Acceptable
		x := 2 // Issue
		fmt.Println(x)
	}
	x = x + 2
	fmt.Println(x)
}

func bar() { //nested test
	x := "test"
	{
		if x == "test" {
			x := "newVariable" // Issue
			fmt.Println(x)
		}
	}
}
