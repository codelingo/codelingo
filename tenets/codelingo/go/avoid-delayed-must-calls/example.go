package main

import (
	"regexp"
	"fmt"
)

func main() {
	foo()
	bar()

	red := regexp.MustCompile("a")
	fmt.Printf(red.String())
}

func foo() {

	var green, _ = regexp.Compile("a")

	yellow, _ := regexp.Compile("a")

	red := regexp.MustCompile("a") // Issue

	fmt.Printf("%s %s %s \n", green.String(), yellow.String(), red.String())
}

func bar() {
	var green, _ = regexp.Compile("a")

	yellow, _ := regexp.Compile("a")

	red := regexp.MustCompile("a") // Issue

	indigo := regexp.MustCompile("a").Match([]byte(`seafood`)) // Issue

	if indigo == true {
		fmt.Printf("%s %s %s \n", green.String(), yellow.String(), red.String())
	}
}

func init() {
	var green, _ = regexp.Compile("a")

	yellow, _ := regexp.Compile("a")

	red := regexp.MustCompile("a") // Acceptable

	fmt.Printf(green.String() + yellow.String() + red.String())
}
