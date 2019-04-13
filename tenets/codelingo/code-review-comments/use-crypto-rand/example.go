//Package main is an example package
package main

import (
	myRand "math/rand"
	// "encoding/base64"
	// "encoding/hex"
	"fmt"
)

// Key is an example function
func Key() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err) // out of randomness, should never happen
	}
	return fmt.Sprintf("%x", buf)
	// or hex.EncodeToString(buf)
	// or base64.StdEncoding.EncodeToString(buf)
}

func main() {
	var key = myRand.Float64()
}
