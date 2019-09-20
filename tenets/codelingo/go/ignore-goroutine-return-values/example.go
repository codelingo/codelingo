package main

import (
	"time"
	"fmt"
)

func main() {
	go A() // Issue
        go B() // Non Issue
	<-time.After(time.Second)

	go func() string { // Issue
	    word := "Hello"
	    return word 
        }()




        go func() { // Non Issue
            fmt.Println("Hello")
        }()

}

func A () string {
	return "str"
}

func B () {
	fmt.Println("Hey")
}
