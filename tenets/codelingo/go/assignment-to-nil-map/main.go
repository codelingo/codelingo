package main

import "time"

func main() {
	now := time.Now().Unix()
	metrics := []string{"example"}

	var accessTimes map[string]int64 // ISSUE
	for _, m := range metrics {
		accessTimes[m] = now
	}
}
