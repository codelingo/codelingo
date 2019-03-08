//Package main is an example package
package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
	a, _ := example() // ISSUE
	fmt.Println(a)
}

func example() (int, error) {
	return 1, errors.New("some error")
}

type a struct{}

func trickyReturnExample() (a, *a, int, *string, string, bool, error) {
	i, _ := example() // ISSUE
	return a{}, nil, i, nil, "hello", true, nil
}
