package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		select {
		case <-time.After(50 * time.Second):
			fmt.Println("Hello, playground")
		}
	}
}

func other() {
	if true {
		for {
			select {
			case <-time.After(50 * time.Second):
				fmt.Println("Hello, playground")
			}
		}
	}
}
