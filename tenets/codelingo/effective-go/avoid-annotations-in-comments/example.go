// Main package for tenet examples
package main

import (
	"fmt"
)

func main() {
	//This is a comment
	fmt.Println("Hello, playground")
	/* This is another comment _this_
	 *
	 */
	fmt.Println("Hello")
	// This comment is fine and should not be picked up
	fmt.Println("Hello")
	// <this> tag should be caught
}
