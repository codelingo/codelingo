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

	slice := []int{0,1,2,3}

	// Issue
	for _, i := range slice {
		go func() {
			fmt.Println(i)
		}()
	}

	// Issue
	for j, i := range slice {
		go func() {
			fmt.Println(j)
		}()
	}

	// Non Issue
	for _, i := range slice {
		go func(i int) {
			fmt.Println(i)
		}(i)
	}

	// Non Issue
	for j, i := range slice {
		go func(j int) {
			fmt.Println(j)
		}(j)
	}

	size := 4

	// Non Issue, the goroutines in this loop use the upper bound of the loop which doesn't change therefore there is no issue
	for i := 0; i < size; i++ {
		go func() {
			fmt.Println(size)
		}()
	}
}
