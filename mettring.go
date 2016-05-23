package mettring

import (
	"time"

	"github.com/ChristianKniep/QNIBCollect/src/fullerite/metric"
)


/* This library provides a RingBuffer based on the attribute .Time
 * It creates buckets of time.Duration and the means to clean up to a given duration period
 * The goal is to buffer objects for a given amount of time in a ringbuffer fashion
 */

// Ring provides the main struct
type Ring struct {
	retention int
	head time.Time
	tail time.Time
	buffer map[int64][]metric.Metric
}

// New returns a new ring
func New(sec int) Ring {
	return Ring{
		retention: sec,
		head: time.Now(),
		tail: time.Now(),
		buffer: make(map[int64][]metric.Metric),
	}
}

/*
Enqueue a value into the Ring buffer.
*/
func (r *Ring) Enqueue(m metric.Metric) {
	// Would be cool if the buckets are not Unix-Epochs but time.Time
	now := m.GetTime()
	ts := now.Unix()
	if r.head.Unix() < ts {
		r.head = now
	}
	_, ok := r.buffer[ts]
	if !ok {
		r.buffer[ts] = make([]metric.Metric, 0)
	}
	r.buffer[ts] = append(r.buffer[ts], m)
}

// Peek returns the slice of a given timestamp
func (r *Ring) Peek(ts int64) ([]metric.Metric, bool) {
	slice, ok := r.buffer[ts]
	if !ok || len(slice) == 0 {
		ret := make([]metric.Metric, 1)
		return ret, false
	}
	return slice, true
}
