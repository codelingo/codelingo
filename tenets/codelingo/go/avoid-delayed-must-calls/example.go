package main

import (
	"regexp"
	"fmt"
)

func main() {
	foo()
	bar()

	red := regexp.MustCompile("a")
	fmt.Println(red.String())
}

func foo() {

	var green, _ = regexp.Compile("a")

	yellow, _ := regexp.Compile("a")

	red := regexp.MustCompile("a") // Issue

	fmt.Println("%s %s %S \n", green.String(), yellow.String(), red.String())
}

func bar() {
	var green, _ = regexp.Compile("a")

	yellow, _ := regexp.Compile("a")

	red := regexp.MustCompile("a") // Issue

	indigo := regexp.MustCompile("a").Match([]byte(`seafood`)) // Issue

	fmt.Println("%s %s %s %s \n", green.String(), yellow.String(), red.String(), indigo)
}

func init() {
	var green, _ = regexp.Compile("a")

	yellow, _ := regexp.Compile("a")

	red := regexp.MustCompile("a") // Acceptable

	fmt.Println(green.String() + yellow.String() + red.String())
}
