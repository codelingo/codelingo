package main

import (
	"time"
	"fmt"
)

func main() {
	go A() // Issue
        go B() // Non Issue
	<-time.After(time.Second)

	go func() string {
	    word := "Hello"
	    return word // Issue
        }()




        go func() {
            fmt.Println("Hello") // Non Issie
        }()

}

func A () string {
	return "str"
}

func B () {
	fmt.Println("Hey")
}
