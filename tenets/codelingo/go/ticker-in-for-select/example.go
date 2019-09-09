package main

import (
	"fmt"
)

func main() {
	for {
		select {
		case <-time.After(50 * time.Second):
			fmt.Println("Hello, playground")
		}
	}

	for _, i := range nums {
		select {
		case <-time.After(50 * time.Second):
			fmt.Println("Hello, playground")
		}
	}
}
