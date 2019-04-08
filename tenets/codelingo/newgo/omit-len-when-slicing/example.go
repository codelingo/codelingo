package main

import (
	"fmt"
)

func main() {
	s := []int{1, 2, 3, 4, 5}
	t := s[1:len(s)]

	fmt.Println(t)
}