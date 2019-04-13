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

func passing() error {
	_, err := example()
	if err != nil {
		return err
	}

	i, err := example()
	if err != nil {
		return err
	}

	_ = i
	return nil
}

func example() (int, error) {
	return 1, errors.New("some error")
}

type a struct{}

func trickyReturnExample() (a, *a, int, *string, string, bool, error) {
	i, _ := example() // ISSUE
	return a{}, nil, i, nil, "hello", true, nil
}

func singleExample() error {
	i, _ := example() // ISSUE
	_ = i
	return nil
}

func nonErrorDiscard() (int, error) {
	_, err := example()
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (b *a) methodExample() (int, error) {
	i, _ := example() // ISSUE
	return i, nil
