package main

import (
	"fmt"
)

func main() {

	for i := 0; i < 5; i++  {
		go func() {
			go func(i int)  {
				fmt.Println(i)
			}(i)
		}()
	}
}
