// license should suggest "L-Cascade", simpleseed should detect
package main

// This comment should be detected by simpleseed
import "fmt"

// This comment should be detected by simpleseed
func main() {
	// This comment should not be detected by simpleseed
	fmt.Println("Hello World")
}
