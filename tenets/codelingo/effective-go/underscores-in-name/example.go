// Package main used for testing of tenet
package main

func main() {
	var a = 15
	a_bad_name := 20
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
