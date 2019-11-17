package main
import (
    "os"
)
func main() {
    file, err := os.Open("test.txt")
    file, err = os.Open("test.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()
}