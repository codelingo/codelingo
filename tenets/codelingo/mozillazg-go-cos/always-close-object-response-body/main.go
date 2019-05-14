package main

import "fmt"

type object int

type client struct {
	Object object
}

func main() {
	var c client
	_, err := c.Object.Put() // ISSUE
	if err != nil {
		panic(err)
	}
}

func (o *object) Put() (int, error) {
	fmt.Println("put")
	return 1, nil
}
