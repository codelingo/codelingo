package main

import "fmt"

type Cat interface {
	Meow()
}

type Tabby struct{}

func (*Tabby) Meow() {
	fmt.Println("meow")
}

func getACat() Cat {
	var myTabby *Tabby = nil

	if false {
		return myTabby
	}
	return myTabby
}

func returnsNil() error {
	var p = nil

	return p
}

func main() {
	if getACat() == nil {
		fmt.Println("Forgot to return a real cat!")
	}
}
