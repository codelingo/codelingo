package main

// Heuristic:
// ‾‾‾‾‾‾‾‾‾‾
// When there is only ONE lock inside the same scope, ensure its name is 'lock'.

import "fmt"

type anyStruct struct {
        lock sync.Mutex
}

func afunc(c *anyStruct) {
        c.lock.Lock()
        defer c.lock.Unlock()
}

