// Package main used for testing of tenet
package main

func main() {
	var a int = 15
	a_bad_name := 20
	var another_bad_name string = "hello"
	aGoodVar := 2
}

func aFunction() {
	fmt.Println("This name is fine")
}

func a_bad_function() {
	fmt.Println("Needs a name change")
}

type Bad_Interface interface {
	SomeMethod() int
}

type BetterInterface interface {
	SomeMethod() int
}
