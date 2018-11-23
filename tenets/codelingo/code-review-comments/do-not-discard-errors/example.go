//Package main is an example package
package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
	a, _ := example()
	fmt.Println(a)
	fmt.Println(err)
}

func example() (int, error) {
	return 1, errors.New("some error")
}
