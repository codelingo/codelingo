package main

import (
	"sync"
)

type Foo struct {
	mu    sync.Mutex
	name  string
	count int
}

// This method is unsafe
func (f *Foo) UnsafeSetName(name string) {
	f.mu.Lock()
	f.name = name
	f.mu.Unlock()
	f.count++ // Issue
}

// This method is safe
func (f *Foo) SafeSetName(name string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.name = name
	f.count++
}

// This method is safe
func (f *Foo) SetName(name string) {
	f.mu.Lock()
	f.name = name
	f.count++ // Accepted
	f.mu.Unlock()
}

func NewFoo() *Foo {
	return &Foo{
		mu:    sync.Mutex{},
		name:  "foo",
		count: 0,
	}
}

func main() {
	wg := sync.WaitGroup{}
	foo := NewFoo()
	wg.Add(5)
	go func() {
		defer wg.Done()
		foo.SetName("bar")
	}()
	go func() {
		defer wg.Done()
		foo.SetName("bar")
	}()
	go func() {
		defer wg.Done()
		foo.SetName("bar")
	}()
	go func() {
		defer wg.Done()
		foo.UnsafeSetName("bar")
	}()
	go func() {
		defer wg.Done()
		foo.UnsafeSetName("bar")
	}()
	wg.Wait()
}
