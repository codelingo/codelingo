package main

import (
	"fmt"
)

type File struct {
	fd   int
	name string
}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func NewFile1(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	f := new(File)
	f.fd = fd     //ISSUE
	f.name = name //ISSUE
	return f
}

func NewFile2(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	f := File{fd, name}
	return &f
}

func main() {
	fmt.Println(NewFile1(5, "Test"))
	fmt.Println(NewFile2(5, "Test"))

	p1 := Person{FirstName: "John", LastName: "Snow", Age: 45}
	fmt.Println(p1)

	p2 := new(Person)
	p2.FirstName = "Alice" //ISSUE
	p2.LastName = "Green"  //ISSUE
	p2.Age = 40            //ISSUE
	fmt.Println(p2)

}
