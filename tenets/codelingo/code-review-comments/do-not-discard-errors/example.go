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
	fmt.Println(err)
}

func example() (int, error) {
	return 1, errors.New("some error")
}

func trickyReturnExample() (int, *string, string, bool, error) {
	i, _ := example() // ISSUE
	return i, nil, "hello", true, nil
}
