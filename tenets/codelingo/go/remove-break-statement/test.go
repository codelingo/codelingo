package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Print("Go runs on ")
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
		break
	case "linux":
		fmt.Println("Linux.")
		break
	default:
		fmt.Printf("%s.\n", os)
		break
	}
}
