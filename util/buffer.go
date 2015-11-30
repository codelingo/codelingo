package util

import "sync/atomic"

type buffer struct {
	// number of allowed tenets
	limit int64

	// number of open tenets
	count int64

	fullc chan struct{}
	roomc chan struct{}
	killc chan struct{}
}

// TODO(waigani) limit must be > 0
func NewBuffer(limit int, killc chan struct{}) *buffer {
	b := &buffer{
		limit: int64(limit),
		fullc: make(chan struct{}),
		roomc: make(chan struct{}),
		count: int64(1),
		killc: killc,
	}
	return b
}

func (b *buffer) Add(i int) {
	atomic.AddInt64(&b.count, int64(i))

	if b.count >= b.limit {
		b.fullc = make(chan struct{})
		close(b.fullc)
		return
	}
	b.roomc = make(chan struct{})
	close(b.roomc)
	return
}

// WaitFull will block until it's full, at which point the returned chan will be
// closed. Future calls to Full will
func (b *buffer) WaitFull() {

	b.Add(0)
	select {
	case <-b.killc:
	case <-b.fullc:
		b.fullc = nil
	}
	return
}

func (b *buffer) Count() int64 {
	return atomic.LoadInt64(&b.count)
}

// Waits until there is room in the buffer. Calculated when buffer is added to
func (b *buffer) WaitRoom() {

	b.Add(0)
	select {
	case <-b.killc:
	case <-b.roomc:
		b.roomc = nil
	}
	return
}
