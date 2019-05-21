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
func test1() *os.File {
	f, err := os.Open("/tmp/test.txt")
	check(err)
	return f
}

func test2() {
	f, err := os.Open("/tmp/test.txt") //ISSUE
	check(err)
	b := make([]byte, 5)
	n, err := f.Read(b)
	check(err)
	fmt.Printf("%d bytes: %s\n", n, string(b))
	//f.Close()
}

func test3() (f *os.File) {
	f, err := os.Open("/tmp/test.txt")
	check(err)
	return
}

func main() {

	f1, err := os.Open("/tmp/test.txt") //ISSUE
	check(err)
	b1 := make([]byte, 5)
	n1, err := f1.Read(b1)
	check(err)
	fmt.Printf("%d bytes: %s\n", n1, string(b1))
	//f1.Close()

	f2, err := os.Open("/tmp/test.txt")
	check(err)
	b2 := make([]byte, 5)
	n2, err := f2.Read(b2)
	check(err)
	fmt.Printf("%d bytes: %s\n", n2, string(b2))
	f2.Close()

	f3, err := os.OpenFile("/tmp/test.txt", os.O_RDWR, 0644) //ISSUE
	check(err)
	b3 := make([]byte, 5)
	n3, err := f3.Read(b3)
	check(err)
	fmt.Printf("%d bytes: %s\n", n3, string(b3))
	//f3.Close()

	f4, err := os.Open("/tmp/test.txt")
	check(err)
	defer f4.Close()
	b4 := make([]byte, 5)
	n4, err := f4.Read(b4)
	check(err)
	fmt.Printf("%d bytes: %s\n", n4, string(b4))

}
