package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
  privateUsedMethod();
}

func privateUsedMethod() {
  fmt.Println("This private method IS used");
}

func privateUnusedMethod() {
  fmt.Println("This private method is not used");
}
