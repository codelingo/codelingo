package main

import "fmt"
import "math"

func main() {
	var f float64
	var f32 float32
	var f64 float64
	fmt.Printf("%t\n", !math.IsNaN(float64(f)))
	fmt.Printf("%t\n", !math.IsNaN(float64(f32)))
	fmt.Printf("%t\n", !math.IsNaN(float64(f64)))
}
