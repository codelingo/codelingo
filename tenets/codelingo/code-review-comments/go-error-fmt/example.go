// Package main is an example package
package main

import "fmt"

func main() {
	fmt.Errorf("something bad.")
	fmt.Errorf("this is okay")
	fmt.Errorf("This is not okay")
	fmt.Errorf("THIS is okay")
	fmt.Errorf("API call thing.")
	fmt.Errorf("THIS is not okay.")
}
