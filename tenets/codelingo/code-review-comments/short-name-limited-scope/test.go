package main

import "fmt"

func main() {

	var a = "initial"
	fmt.Println(a)

	var b, c int = 1, 2
	fmt.Println(b, c)

	var boolVar = true //ISSUE
	fmt.Println(boolVar)

	var intVar int //ISSUE
	fmt.Println(intVar)

	f1 := "short" //ISSUE
	fmt.Println(f1)
}
