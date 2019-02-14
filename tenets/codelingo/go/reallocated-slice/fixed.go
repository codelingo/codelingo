// Copyright 2017 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package watcher

import (
	"fmt"
)

// HubWatcher listents to events from the hub and passes them on to the registered
// watchers.
type HubWatcher struct {
	syncEvents, requestEvents []event
}

func (w *HubWatcher) flush() {
	// syncEvents are stored first in first out.
	// syncEvents may grow during the looping here if new
	// watch events come in while we are notifying other watchers.
	for i := 0; i < len(w.syncEvents); i++ {
		// We need to reget the address value each time through the loop
		// as the slice may be reallocated.
		for e := &w.syncEvents[i]; e.ch != nil; e = &w.syncEvents[i] {
			w.logger.Tracef("syncEvents: e.ch=%v len(%d), cap(%d)", e.ch, len(w.syncEvents), cap(w.syncEvents))
			select {
			case <-w.tomb.Dying():
				return
			case req := <-w.request:
				w.handle(req)
				continue
			case change := <-w.changes:
				w.queueChange(change)
				continue
			case e.ch <- Change{e.key.c, e.key.id, e.revno}:
				w.logger.Tracef("e.ch=%v has been notified", e.ch)
			}
			break
		}
	}
	w.syncEvents = w.syncEvents[:0]
	w.logger.Tracef("syncEvents: len(%d), cap(%d)", len(w.syncEvents), cap(w.syncEvents))

	// requestEvents are stored oldest first, and
	// may grow during the loop.
	for i := 0; i < len(w.requestEvents); i++ {
		// We need to reget the address value each time through the loop
		// as the slice may be reallocated.
		for e := &w.requestEvents[i]; e.ch != nil; e = &w.requestEvents[i] {
			select {
			case <-w.tomb.Dying():
				return
			case req := <-w.request:
				w.handle(req)
				continue
			case change := <-w.changes:
				w.queueChange(change)
				continue
			case e.ch <- Change{e.key.c, e.key.id, e.revno}:
			}
			break
		}
	}
	w.requestEvents = w.requestEvents[:0]
}

// handle deals with requests delivered by the public API
// onto the background watcher goroutine.
func (w *HubWatcher) handle(req interface{}) {
	w.logger.Tracef("got request: %#v", req)
	switch r := req.(type) {
	case reqWatch:
		for _, info := range w.watches[r.key] {
			if info.ch == r.info.ch {
				panic(fmt.Errorf("tried to re-add channel %v for %s", info.ch, r.key))
			}
		}
		if revno, ok := w.current[r.key]; ok && (revno > r.info.revno || revno == -1 && r.info.revno >= 0) {
			r.info.revno = revno
			w.requestEvents = append(w.requestEvents, event{r.info.ch, r.key, revno})
		}
		w.watches[r.key] = append(w.watches[r.key], r.info)
	case reqUnwatch:
		watches := w.watches[r.key]
		removed := false
		for i, info := range watches {
			if info.ch == r.ch {
				watches[i] = watches[len(watches)-1]
				w.watches[r.key] = watches[:len(watches)-1]
				removed = true
				break
			}
		}
		if !removed {
			panic(fmt.Errorf("tried to remove missing channel %v for %s", r.ch, r.key))
		}
		for i := range w.requestEvents {
			e := &w.requestEvents[i]
			if r.key.match(e.key) && e.ch == r.ch {
				e.ch = nil
			}
		}
		for i := range w.syncEvents {
			e := &w.syncEvents[i]
			if r.key.match(e.key) && e.ch == r.ch {
				e.ch = nil
			}
		}
	default:
		panic(fmt.Errorf("unknown request: %T", req))
	}
}

// queueChange queues up the change for the registered watchers.
func (w *HubWatcher) queueChange(change Change) {
	w.logger.Tracef("got change document: %#v", change)
	key := watchKey{change.C, change.Id}
	revno := change.Revno
	w.current[key] = revno

	// Queue notifications for per-collection watches.
	for _, info := range w.watches[watchKey{change.C, nil}] {
		if info.filter != nil && !info.filter(change.Id) {
			continue
		}
		w.syncEvents = append(w.syncEvents, event{info.ch, key, revno})
		w.logger.Tracef("adding collection watch for %v syncEvents: len(%d), cap(%d)", info.ch, len(w.syncEvents), cap(w.syncEvents))
	}

	// Queue notifications for per-document watches.
	infos := w.watches[key]
	for i, info := range infos {
		if revno > info.revno || revno < 0 && info.revno >= 0 {
			infos[i].revno = revno
			w.syncEvents = append(w.syncEvents, event{info.ch, key, revno})
			w.logger.Tracef("adding document watch for %v syncEvents: len(%d), cap(%d)", info.ch, len(w.syncEvents), cap(w.syncEvents))
		}
	}
}
