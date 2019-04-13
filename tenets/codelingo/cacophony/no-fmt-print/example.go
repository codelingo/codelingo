package main

import (
	"fmt"
	"log"
)

func main() {
	correct()
	incorrect()
}

func correct() {
	log.Println("correct called")
	a := 5
	log.Printf("value of a is %d\n", a)
}

func incorrect() {
	fmt.Println("incorrect called") // ISSUE
	a := 7
	fmt.Printf("value of a is %d\n", a) // ISSUE
}
