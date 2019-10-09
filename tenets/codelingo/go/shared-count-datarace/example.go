package main

import (
	"fmt"
	"sync"
)

func main() {

	// Issue - When a loop creates goroutines which needs to use the iterator of the
	// loop the iterator should always be provided as an argument to the function.
	// If this is not done, as in the below case, then we risk having an unexpected
	// value of the iterator as the parent goroutine may have incremented its value
	// already.
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

	// Non Issue, the use of a WaitGroup here means there is no risk of a race condition
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Prinltn(i)
		}()
	}
	wg.Wait()
}
