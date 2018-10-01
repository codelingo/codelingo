package main

import "fmt"

var privateVarA string = "baz"
var privateVarB string = "baz"

func foo() {
	fmt.Printf("%s\n", privateVarA)
}
