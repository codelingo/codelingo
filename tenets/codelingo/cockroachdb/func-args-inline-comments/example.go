package main

import (
	"fmt"
)

func main() {
	// Well commented
	monitor := mon.MakeMonitor(
		"in-mem temp storage",
		mon.MemoryResource,
		nil,             /* curCount */
		nil,             /* maxHist */
		1024*1024,       /* increment */
		maxSizeBytes/10, /* noteworthy */
	)

	// Missing comments
	monitor := mon.MakeMonitor(
		"in-mem temp storage",
		mon.MemoryResource,
		nil,
		nil,
		1024*1024,
		maxSizeBytes/10,
	)
}
