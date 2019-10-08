package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
}

func f(c chan int) {} // ISSUE

func g(c <-chan int) {}

func h(c chan<- int) {}

func i() chan int { // ISSUE
	return nil
}

func j() <-chan int {
	return nil
}

func k() chan<- int {
	return nil
}
