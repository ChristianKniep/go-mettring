package mettring

import (
	//"fmt"
	"time"
)

// Item provides the interface that expects Time()
type Item interface {
	Time() time.Time
}
/* This library provides a RingBuffer based on the attribute .Time
 * It creates buckets of time.Duration and the means to clean up to a given duration period
 * The goal is to buffer objects for a given amount of time in a ringbuffer fashion
 */

// Ring provides the main struct
type Ring struct {
	retention int
	head time.Time
	tail time.Time
	buffer map[int64][]interface{}
}

// New returns a new ring
func New(sec int) Ring {
	return Ring{
		retention: sec,
		head: time.Now(),
		tail: time.Now(),
	}
}

/*
Enqueue a value into the Ring buffer.
*/
func (r *Ring) Enqueue(i Item) {
	// Would be cool if the buckets are not Unix-Epochs but time.Time
	now := i.Time()
	ts := now.Unix()
	if r.head.Unix() < ts {
		r.head = now
	}
	_, ok := r.buffer[ts]
	if !ok {
		r.buffer[ts] = make([]interface{}, 20)
	}
	r.buffer[ts] = append(r.buffer[ts], i)
}

// Peek returns the slice of a given timestamp
func (r *Ring) Peek(ts int64) ([]interface{}, bool) {
	slice, _ := r.buffer[ts]
	if !ok || len(slice) == 0 {
		ret := make([]interface{}, 2)
		ret = append(ret, NewTest("nil", "nil"))
		return ret, false
	}
	return slice, true
}
