package main

import "fmt"
import "testing"

type Cat interface {
	Meow()
}

type Tabby struct{}

func (*Tabby) Meow() { fmt.Println("meow") }

// Must return something that is like 'Cat'. A 'Tabby' is like 'Cat' because it has a 'Meow' function
func GetACat() Cat {
	var myTabby *Tabby = nil
	// Oops, we forgot to set myTabby to a real value
	return myTabby
}

func TestGetACat(t *testing.T) {
	// This test does not do what was intended. It's never nil because it's a pointer to a pointer
	if GetACat() == nil {
		t.Errorf("Forgot to return a real cat!")
	}
}

// I am only testing the above circumstances, precisely, for the the moment

func main() {
}

// Heuristic:
// ‾‾‾‾‾‾‾‾‾‾
// look for 'varname == nil' test where varname is a pointer (of type interface pointer?) to a pointer to nil.
