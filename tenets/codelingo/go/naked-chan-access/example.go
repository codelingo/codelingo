package main

import (
	"time"
	"fmt"
)

func inc(i int, c chan int) {
    c <- i + 1
}

func main() {

	c := make(chan int)

	go inc(0, c)

	x := <- c // Issue 
	fmt.Println(x)

	go inc(4, c)

	select {
	case i := <- c: // Safe
		fmt.Println(i)
	case <- time.After(1 * time.Second):
		fmt.Println("Timed out")
	}

	go inc(3, c)

	select {
	case x = <- c: // Safe
		fmt.Println(x)
	case <- time.After(1 * time.Second):
		fmt.Println("Timed out")
	}
}
