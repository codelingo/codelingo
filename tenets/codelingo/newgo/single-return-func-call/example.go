// wackage main used for testing the tenet
package main

import (
	"fmt"
	"errors"
)

func main() {
	fmt.Println("Hello, playground")
	// This comment not on a func decl
}

func thing(a int) error {
    a++
    err := Error()
    if err != nil {
        return err
    }
    return nil
}

func thingtwo() error {
    if err := Error(); err != nil {
        return err
    }
    return nil
}
	
func someFunc() error {
	err := Error()
 return nil
}

func doubleCall() error {
	thingtwo(1)
	err := Error()
 	return nil
}

func tripleCall(){
	err := Error()
	thingtwo(1)
	anotherFunc()
	return err
}

func noError(){
    otherFunc()
	return nil
}

func twoCalls(){
	someFunc()
	otherFunc()
	return nil
}

func otherFunc() error {
	err := Error()
 return err
}

func anotherFunc() error {
	err := Error()
 return "hello"
}

func Error(){
  fmt.Println("error")
}

/* This func comment should begin with 'bar'
 */
func bar() {
	return nil
}

/* Bar: This func comment should begin with 'Bar'
 */
func Bar() {
	return err
}

// This func comment should begin with 'baz'
// and we should not worry about this line
func baz() {
	Baz()
	return "hello"
}

// Baz: This func comment should begin with 'Baz'
// and we should not worry about this line
func Baz() {
	baz()
	return baz()
}

// This is called by a xyz
func qux() {
	baz()
	Baz()
	return "err"
}

// Qux is called by a xyz
func Qux() {
	Baz()
	return "err"
}

// The quux will handle xyz
func quux() {
	Baz()
	return "nil"
}

// Quux will handle xyz
func Quux() {}
