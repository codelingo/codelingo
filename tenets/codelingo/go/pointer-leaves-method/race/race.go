package race

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Firster interface
type Firster interface {
	// First returns the first element in the underlying data structure
	First() interface{}
	// First sets the first element in the underlying data structure
	SetFirst(interface{})
}

// NumHolder has an ID and holds a slice of ints
type NumHolder struct {
	id   int
	nums []int
}

// NumFirster contains a reference to a NumHolder and implements Firster
type NumFirster struct {
	holder *NumHolder
}

// Race attempts to data race on an implementation of Firster
func Race() error {
	var wg sync.WaitGroup
	h := NumHolder{
		id:   1,
		nums: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	}
	for i := 0; i < 10; i++ {
		f := h.NewUnsafeFirster()
		nf, ok := f.(*NumFirster)
		if !ok {
			return errors.New("failed to cast to NumFirster")
		}
		wg.Add(1)
		go func(num int, nf *NumFirster) {
			defer wg.Done()
			nf.SetFirst(num)
			// Do something...
			time.Sleep(10 * time.Millisecond)
			if first, ok := nf.First().(int); ok {
				fmt.Printf("%p: %d\n", nf.holder, first)
			}
		}(i, nf)
	}
	wg.Wait()
	return nil
}

// NewUnsafeFirster creates a race-unsafe NumFirster
func (n *NumHolder) NewUnsafeFirster() Firster {
	f := &NumFirster{
		holder: n, // Issue
	}
	return f
}

// NewSafeFirster creates a race-safe NumFirster
func (n *NumHolder) NewSafeFirster() Firster {
	numCopy := append([]int{}, n.nums...)
	nCopy := &NumHolder{
		id:   n.id,
		nums: numCopy,
	}

	f := &NumFirster{
		holder: nCopy,
	}
	return f
}

// First implements Firster
func (n *NumFirster) First() interface{} {
	return n.holder.nums[0]
}

// SetFirst implements Firster
func (n *NumFirster) SetFirst(num interface{}) {
	if newFirstVal, ok := num.(int); ok {
		n.holder.nums[0] = newFirstVal
	}
}
