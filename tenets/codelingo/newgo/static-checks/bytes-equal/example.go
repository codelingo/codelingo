package main

import (
	"fmt"
)

func main() {

	s := []byte("Hello")
	t := []byte("Hell")

	fmt.Println(bytes.Compare(s, t) == 0)

	fmt.Println(bytes.Compare("word", "wo") == 0)

}