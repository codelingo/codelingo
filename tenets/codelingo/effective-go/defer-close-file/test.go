package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	f1, err := os.Open("/tmp/test.txt") //ISSUE
	check(err)
	//defer f1.Close()
	b1 := make([]byte, 5)
	n1, err := f1.Read(b1)
	check(err)
	fmt.Printf("%d bytes: %s\n", n1, string(b1))

	f2, err := os.Open("/tmp/test.txt") 
	check(err)
	defer f2.Close()
	b2 := make([]byte, 5)
	n2, err := f2.Read(b2)
	check(err)
	fmt.Printf("%d bytes: %s\n", n2, string(b2))
	f2.Close()

	f3, err := os.OpenFile("/tmp/test.txt", os.O_RDWR, 0644) //ISSUE
	check(err)
	//defer f3.Close()
	b3 := make([]byte, 5)
	n3, err := f3.Read(b3)
	check(err)
	fmt.Printf("%d bytes: %s\n", n3, string(b3))
      
}
