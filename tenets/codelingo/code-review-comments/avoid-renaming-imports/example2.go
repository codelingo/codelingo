// Package main is an example package
package main

import (
	badfmt "fmt"
	rando "math/rand"
	_ "os"
)

func main() {
	badfmt.Println("Hello, playground", rando.New())
}
