package main

import (
	"fmt"
	"os"
)

func main() {
	foo()
	bar()
	baz()
}

func baz() {
	file, err := os.Open("test.txt")
	if err != nil {
		return 
	}
	defer file.Close()

	fileList := []string{"test1.txt"}
	for _, f:= range fileList {
		fileN, err := os.Open(f) // Acceptable
		if err != nil {
			return 
		}
		defer fileN.Close()
	}
}

func foo() {
	x := 1	// Acceptable
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
		if true {
			x := "newVariable" // Issue
			fmt.Println(x)
		}
	}
	fmt.Println(x)
}