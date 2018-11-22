package main

import (
	"fmt"
)

/*
  The problem with this comment is that it exceeds 80 characters and needs to be wrapped
*/
func main() {
	fmt.Println(`Hello, world! This function call has more than 100 chars 
        but it has been wrapped, so we don't want to find this`)
	fmt.Println("Hello, world! The issue with this line of code is that it is over 100 chars and needs to be wrapped")
	fmt.Println("This line is fine")
}
