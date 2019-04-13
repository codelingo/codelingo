package main

import "context"

// Example is an example struct
type Example struct {
	Name string
	Ctx  context.Context
}

func main() {
	var ex Example
	ex.Ctx = context.Background()
	ex.Name = "Example"
}
