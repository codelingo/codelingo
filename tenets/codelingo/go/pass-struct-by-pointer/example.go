package main

import (
	"fmt"
)

type Cat struct {
	Name string
	age int
}

type dog struct {
	name string
	age int
}

func getCatName(c *Cat) {
	fmt.Println(c.Name) // Issue, accessing struct with exported field(s) by reference
}

func getCatAge(c Cat) { // Non Issue, passing struct by value
	fmt.Println(c.Age)
}

func getDogName(d *dog) { // Non Issue, acessing struct with no exported fields
	fmt.Println(d.Name)
}

func getDogAge(d dog) { // Non Issue, passing struct by value with no exported fields
	fmt.Println(d.age)
}

func (c *Cat) getName() { // Issue, accesing struct with exported field(s) by reference
	fmt.Pritnln(c.Name)
}

func (c Cat) getAge() { // Non Issue, passing struct by value
	fmt.Println(c.age)
}

func main() {
	fmt.Println("Hey")

}
