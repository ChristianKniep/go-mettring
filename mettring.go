package mettring

import (
	"time"
	//"fmt"
	"sort"
	"strconv"

	"github.com/ChristianKniep/QNIBCollect/src/fullerite/metric"
)

/* This library provides a RingBuffer based on the attribute .Time
 * It creates buckets of time.Duration and the means to clean up to a given duration period
 * The goal is to buffer objects for a given amount of time in a ringbuffer fashion
 */

// Itoa64 returns string from int64
func Itoa64(i int64) string {
	return strconv.FormatInt(i, 10)
}

// Ring provides the main struct
type Ring struct {
	retention int
	head time.Time
	tail time.Time
	count int
	buffer map[string][]metric.Metric
}

// New returns a new ring
func New(ms int) Ring {
	return Ring{
		retention: ms,
		head: time.Now(),
		tail: time.Now(),
		count: 0,
		buffer: make(map[string][]metric.Metric,1),
	}
}

/*
Enqueue a value into the Ring buffer.
*/
func (r *Ring) Enqueue(m metric.Metric) {
	// Would be cool if the buckets are not Unix-Epochs but time.Time
	now := m.GetTime()
	ts := now.UnixNano()
	if r.head.UnixNano() < ts {
		r.head = now
	}
	sts := Itoa64(ts)
	_, ok := r.buffer[sts]
	if !ok {
		r.buffer[sts] = []metric.Metric{}
	}
	r.buffer[sts] = append(r.buffer[sts], m)
	r.count++
}

// Peek returns the slice of a given timestamp
func (r *Ring) Peek(ts int64) ([]metric.Metric, bool) {
	sts := Itoa64(ts)
	slice, ok := r.buffer[sts]
	if !ok || len(slice) == 0 {
		ret := []metric.Metric{}
		return ret, false
	}
	return slice, true
}

// Values returns the buffer
func (r *Ring) Values() ([]metric.Metric, bool) {

	if len(r.buffer) == 0 {
		ret := []metric.Metric{}
		return ret, false
	}
	var keys []string
	for k := range r.buffer {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ret := []metric.Metric{}
	for _, key := range keys {
		for _, m := range r.buffer[key] {
			ret = append(ret, m)
		}
	}
	return ret, true
}

// TidyUp cleans outdated entries
func (r *Ring) TidyUp() (int, bool) {
	var kicked int
	var keys []string
	for k := range r.buffer {
    keys = append(keys, k)
	}
	sort.Strings(keys)
	old := time.Now().UnixNano() - int64(r.retention*1000000)
	for _, k := range keys {
		ms := r.buffer[k]
		for i, m := range ms {
			if old > m.GetTime().UnixNano() {
				//fmt.Printf("Remove '%s' as old '%d' is more recent then '%d' (%dsec diff)\n", m.Name, old, m.GetTime().UnixNano(), (m.GetTime().UnixNano()-old)/1000000000)
				kicked++
				r.buffer[k] = append(ms[:i], ms[i+1:]...)
			} else {
				//fmt.Printf("Kept '%s' as old '%d' is older then '%d' (%d)\n", m.Name, old, m.GetTime().UnixNano(), m.GetTime().UnixNano()-old)
			}
		}
	}
	return kicked, true
}

// Filter returns the buffer with a filter applied
func (r *Ring) Filter(f metric.Filter) ([]metric.Metric, bool) {

	if len(r.buffer) == 0 {
		ret := []metric.Metric{}
		return ret, false
	}
	var keys []string
	for k := range r.buffer {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	ret := []metric.Metric{}
	for _, key := range keys {
		for _, m := range r.buffer[key] {
			if ! m.IsFiltered(f) {
				ret = append(ret, m)
			}
		}
	}
	return ret, true
}
