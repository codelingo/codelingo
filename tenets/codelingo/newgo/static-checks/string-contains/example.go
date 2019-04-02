package main

import (
	"fmt"
)

func main() {

	s := "Hello"
	t := "Hell"

	fmt.Println(strings.Index(s, t) != -1)

	fmt.Println(strings.Index("word", "wo") != -1)

}
