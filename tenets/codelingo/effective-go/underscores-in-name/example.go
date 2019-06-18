// Package main used for testing of tenet
package main

import "fmt"

func main() {
	var a = 15
	a_bad_name := 20
	HashSHA512_384 := func() {} // A reasonable exception
	var another_bad_name = "hello"
	aGoodVar := 2
}

func aFunction() {
	fmt.Println("This name is fine")
}

func a_bad_function() {
	fmt.Println("Needs a name change")
}

// Bad_Interface used as example
type Bad_Interface interface {
	SomeMethod() int
}

// BetterInterface used as example
type BetterInterface interface {
	SomeMethod() int
}
