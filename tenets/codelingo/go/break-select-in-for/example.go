package main

import "fmt"

func main() {
	// We generally write break statements to break out of repetitious things like for loops
	fmt.Println("Expecting 0, 1:")
	for i := range []int{0, 1, 2} {
		fmt.Println(i)
		if i == 1 {
			break // Non-issue
		}
	}

	// and don't often break from selects, so it's easy to forget that it's possible
	fmt.Println("Expecting 0:")
	c1 := intChan()
	select {
	case i := <-c1:
		fmt.Println(i)
		break // Non-issue
	}

	// We understand that breaking from nested loops only exits the innermost
	fmt.Println("Expecting {0 ,0}, {1, 0}, {2, 0}:")
	for i := range []int{0, 1, 2} {
		for j := range []int{0, 1, 2} {
			fmt.Println(i, j)
			if j == 0 {
				break // Non-issue
			}
		}
	}

	// Unless, of course, we use a label
	fmt.Println("Expecting {0 ,0}")
l:
	for i := range []int{0, 1, 2} {
		for j := range []int{0, 1, 2} {
			fmt.Println(i, j)
			if j == 0 {
				break l // Non-issue
			}
		}
	}

	// Which works perfectly fine to break a select insisde a for loop
	fmt.Println("Expecting 0, 1, 2:")
	c2 := intChan()
m:
	for {
		select {
		case i, ok := <-c2:
			fmt.Println(i)
			if !ok {
				fmt.Println("But actually we get infinite 0s as well")
				break m // Non-issue
			}
		}
	}

	// The trouble is, we often use selects inside for loops and intend to break the for loop
	// forgetting that we need a label lest we only break the select statement
	fmt.Println("Expecting 0, 1, 2:")
	c3 := intChan()
	breakCountDown := 10
	for {
		select {
		case i, ok := <-c3:
			fmt.Println(i)
			if !ok {
				fmt.Println("But actually we get infinite 0s as well")
				breakCountDown--
				if breakCountDown <= 0 {
					panic("escape infinite loop")
				}
				break // ISSUE
			}
		}
	}
}

func intChan() <-chan int {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	return ch
}
