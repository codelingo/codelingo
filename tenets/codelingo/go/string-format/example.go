package main

import "fmt"

func main() {
    errors.Errorf("argument \"%s\" has conflicting definitions")
    errors.Errorf("argument \"%s\" has conflicting definitions", varName)
	errors.Errorf("argument \"%s%s%s\" has conflicting definitions")
	return errors.Errorf("argument \"%s%s%s\" has conflicting definitions", avar1, avar2, avar3)
}

func anotherFunc() {
	fmt.Printf("This ones okay %s", x)
	fmt.Printf("This one is missing %s")
	fmt.Printf("And this one %s%s", y)
	return 	errors.Errorf("argument \"%s%s%s\" has conflicting definitions", avar1, avar2)

}
