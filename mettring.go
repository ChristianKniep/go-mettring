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
	head      time.Time
	tail      time.Time
	count     int
	buffer    map[string]map[string][]metric.Metric
	aggregate map[string]metric.Metric
}

// New returns a new ring
func New(ms int) Ring {
	return Ring{
		retention: ms,
		head:      time.Now(),
		tail:      time.Now(),
		count:     0,
		buffer:    make(map[string]map[string][]metric.Metric),
		aggregate: make(map[string]metric.Metric),
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
		r.buffer[sts] = make(map[string][]metric.Metric)
	}
	id := m.GetFilterID()
	_, ok = r.buffer[sts][id]
	if !ok {
		r.buffer[sts][id] = []metric.Metric{}
	}
	r.buffer[sts][id] = append(r.buffer[sts][id], m)
	r.count++
}

// Peek returns the slice of a given timestamp
func (r *Ring) Peek(ts int64) ([]metric.Metric, bool) {
	sts := Itoa64(ts)
	slice := []metric.Metric{}
	for _, ms := range r.buffer[sts] {
		slice = append(slice, ms...)
	}
	if len(slice) == 0 {
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
		for _, ids := range r.buffer[key] {
			ret = append(ret, ids...)
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
		ids := r.buffer[k]
		for id, ms := range ids {
			for i, m := range ms {
				if old > m.GetTime().UnixNano() {
					//fmt.Printf("Remove '%s' as old '%d' is more recent then '%d' (%dsec diff)\n", m.Name, old, m.GetTime().UnixNano(), (m.GetTime().UnixNano()-old)/1000000000)
					kicked++
					r.buffer[k][id] = append(ms[:i], ms[i+1:]...)
				} else {
					//fmt.Printf("Kept '%s' as old '%d' is older then '%d' (%d)\n", m.Name, old, m.GetTime().UnixNano(), m.GetTime().UnixNano()-old)
				}
			}
		}
	}
	return kicked, true
}

// Filter returns the buffer w/o elements that match the filter
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
		ids := r.buffer[key]
		for _, ms := range ids {
			for _, m := range ms {
				if !m.IsFiltered(f) {
					ret = append(ret, m)
				}
			}
		}
	}
	return ret, true
}

// Match returns the buffer only with elements that match the filter
func (r *Ring) Match(f metric.Filter) ([]metric.Metric, bool) {

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
		ids := r.buffer[key]
		for _, ms := range ids {
			for _, m := range ms {
				if m.IsFiltered(f) {
					ret = append(ret, m)
				}
			}
		}
	}
	return ret, true
}

// GetAggregate returns the aggregated buffer
func (r *Ring) GetAggregate() ([]metric.Metric, bool) {

	if len(r.aggregate) == 0 {
		ret := []metric.Metric{}
		return ret, false
	}
	var ids []string
	for i := range r.aggregate {
		ids = append(ids, i)
	}
	ret := []metric.Metric{}
	for _, i := range ids {
			ret = append(ret, r.aggregate[i])
	}
	return ret, true
}

// AggregateBuffer iterates over the buffer and aggregates the list of IDs into a single metric
func (r *Ring) AggregateBuffer() bool {
	var keys []string
	for k := range r.buffer {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	r.aggregate = make(map[string]metric.Metric)
	cacheSlice := make(map[string][]metric.Metric)
	var agg, cache float64
	for _, key := range keys {
		for id, ms := range r.buffer[key] {
			_, ok := cacheSlice[id]
			if !ok {
				cacheSlice[id] = []metric.Metric{}
			}
			for _, m := range ms {
				cacheSlice[id] = append(cacheSlice[id], m)
			}
		}
	}
	for id, ms := range cacheSlice {
		var m metric.Metric
		cache = 0.0
		for _, m = range ms {
			cache = cache + m.Value
		}
		agg = cache / float64(len(ms))
		aggMetric := metric.NewExt(m.Name, m.MetricType, agg, m.Dimensions, m.GetTime(), m.Buffered)
		r.aggregate[id] = aggMetric
	}
	return true
}
