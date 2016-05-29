package mettring

import (
	"fmt"
	"testing"
	"time"

	"github.com/ChristianKniep/QNIBCollect/src/fullerite/metric"
	"github.com/stretchr/testify/assert"
)

func TestNewMettring(t *testing.T) {
	mr := New(300)
	assert.Equal(t, mr.retention, 300)
}

func TestMettringEnque(t *testing.T) {
	mr := New(300)
	m := metric.New("TestMetric")
	now := time.Now()
	m.SetTime(now)
	mr.Enqueue(m)
	item, ok := mr.Peek(now.UnixNano())
	assert.True(t, ok, "nothing in the slice")
	assert.Equal(t, m, item[0])
}

func TestMettringValues(t *testing.T) {
	mr := New(200)
	_, ok := mr.Values()
	exp := []metric.Metric{}
	var m metric.Metric
	assert.False(t, ok, "Values() returned non-empty list")
	for i := 0; i < 5; i++ {
		m = metric.New(fmt.Sprintf("m%d", i))
		exp = append(exp, m)
		mr.Enqueue(m)
		time.Sleep(100 * time.Millisecond)
	}
	slice, ok := mr.Values()
	assert.True(t, ok, "Values() returned empty list")
	assert.Equal(t, exp, slice)
}

func TestMettringTidyUp(t *testing.T) {
	mr := New(400) // 400ms retention time
	_, ok := mr.Values()
	exp := []metric.Metric{}
	var m metric.Metric
	assert.False(t, ok, "Values() returned non-empty list")
	// putting 8 items with 125ms distance into it
	for i := 0; i < 8; i++ {
		m = metric.New(fmt.Sprintf("m%d", i))
		exp = append(exp, m)
		mr.Enqueue(m)
		time.Sleep(125 * time.Millisecond)
	}
	// TidyUp should kick out 1
	kicked, ok := mr.TidyUp()
	expKick := 5
	assert.Equal(t, expKick, kicked, fmt.Sprintf("Would have expected to kick %d", expKick))
	assert.True(t, ok, "TidyUp returns false")
	slice, ok := mr.Values()
	assert.True(t, ok, "Values() returned empty list")
	assert.Equal(t, exp[expKick:], slice)
}

func TestMettringFilter(t *testing.T) {
	mr := New(200)
	_, ok := mr.Values()
	exp := []metric.Metric{}
	var m metric.Metric
	assert.False(t, ok, "Values() returned non-empty list")
	for i := 0; i < 5; i++ {
		m = metric.New(fmt.Sprintf("m%d", i))
		exp = append(exp, m)
		mr.Enqueue(m)
		time.Sleep(10 * time.Millisecond)
	}
	d := map[string]string{}
	f := metric.NewFilter("m.*", "gauge", d)
	slice, ok := mr.Filter(f)
	assert.Equal(t, []metric.Metric{}, slice)
	f = metric.NewFilter("m.*", "FAIL_TYPE", d)
	slice, ok = mr.Filter(f)
	assert.Equal(t, exp, slice)
	f = metric.NewFilter("m0", "gauge", d)
	slice, ok = mr.Filter(f)
	assert.Equal(t, exp[1:], slice)
}
