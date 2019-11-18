package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buff := make([]byte, 6)

	n, err := file.Read(buff)
	fmt.Println("read", n, "data:", string(buff))
	if err != nil {
		panic(err)
	}

	n, err = file.Read(buff)
	if err != nil { // Issue
		panic(err)
	}

	fmt.Println("read", n, "data:", string(buff))

}
