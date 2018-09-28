package main

import "fmt"

func main() {
	correct()
	incorrect()
}

func correct() {
	fmt.Println("hello world")
	fmt.Printf("hello %s\n", "world")

	var str string
	str = fmt.Sprintln("hello %s", "world")
	fmt.Print(str)
	str = fmt.Sprintf("4 + 6 = %d\n", 4+6)
	fmt.Print(str)
}

func incorrect() {
	fmt.Println("hello %s", "world")

	var str string
	str = fmt.Sprint("hello %s\n", "world")
	fmt.Print(str)
	str = fmt.Sprint("4 + 6 = %d\n", 4+6)
	fmt.Print(str)
}
