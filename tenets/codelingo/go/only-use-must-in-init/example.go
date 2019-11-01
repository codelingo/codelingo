package main

import (
	"regexp"
	"fmt"
)

func main() {
	foo()
	bar()
	init()

}

func foo() {

	var green, _ = regexp.Compile("a")

	yelllow, _ := regexp.Compile("a")

	red := regexp.MustCompile("a") // Issue

	fmt.Println(green + yellow + red)
}

func bar() {
	var green, _ = regexp.Compile("a")

	yelllow, _ := regexp.Compile("a")

	red := regexp.MustCompile("a") // Issue

	fmt.Println(green + yellow + red)
}

func init() {
	var green, _ = regexp.Compile("a")

	yelllow, _ := regexp.Compile("a")

	red := regexp.MustCompile("a")

	fmt.Println(green + yellow + red)
}
