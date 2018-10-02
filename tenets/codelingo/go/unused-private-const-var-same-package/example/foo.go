package main

import "fmt"

var ExportedVar string = "baz"
var unexportedVar string = "baz"

const ExportedConst string = "baz"
const unexportedConst string = "baz"

type unexportedStructA struct {
	m map[string]bool
}

func unexportedFunc() {
	var b unexportedStructA
	b.m = map[string]bool{
		"exhibit B": true,
	}

	fmt.Printf("%t\n", b.m["exhibit B"])
	fmt.Printf("%s\n", privateVarB)
}

var privateVarA string = "baz"
var privateVarB string = "baz"
var privateVarC string = "baz"

func foo() {
    a := unexportedVar
    fmt.Printf("%s\n", a)

	fmt.Printf("%s\n", privateVarA)
}
