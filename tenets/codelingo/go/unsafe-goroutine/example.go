package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
}
func some() {
 for _, item := range list {
    go func(safeItem int) {
        // do something with item 
        safeItem
    }(safeItem)
 }
}

func someElse() {
 for _, unsafeItem := range list {
    go func() {
        // do something with item 
        unsafeItem
    }()
 }
}

func someElseAgain() {
 for i := 0; i < len(unsafeItem); i++ {
    go func() {
        // do something with item 
        unsafeItem[i]
    }()
 }
}