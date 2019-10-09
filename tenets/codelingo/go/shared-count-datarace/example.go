package main

import (
	"fmt"
)

func main() {

	// Issue - Each goroutine inside this for loop is intended to have a unique
	// value of i. However this may not always be the case so the goroutine should
	// be given a copy of I as an argument as in the second for loop.
	for i := 0; i < 5; i++  {
		go func() {
			fmt.Println(i)
		}()
	}

	// Non Issue
	for i := 0; i < 5; i++ {
		go func(i int) {
			fmt.Println(i)
		}(i)
	}
}
