//Package main is an example package
package main

import (
	"fmt"
)

func main() {
	var s = "Hello, world!"
	var i = 42
	var f float32 = 42.24
	fmt.Println("%v", s)
	fmt.Println("%s", s)
	fmt.Println("%v%v", s, i)
	fmt.Println("%v", i)
	fmt.Println("%v is a string", s)
	fmt.Println("%g", f)
	fmt.Println("%v", f)

	fmt.Println("%s%v", s, i)

	fmt.Println("%s%v%d", s, s, i)
}
