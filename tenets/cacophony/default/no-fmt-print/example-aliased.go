package main

import (
	f "fmt"
	l "log"
)

func main() {
	correctAliased()
	incorrectAliased()
}

func correctAliased() {
	l.Println("correct called")
	a := 5
	l.Printf("value of a is %d\n", a)
}

func incorrectAliased() {
	f.Println("incorrect called") // ISSUE - TODO: this won't be found yet
	a := 7
	f.Printf("value of a is %d\n", a) // ISSUE - TODO: this won't be found yet
}
