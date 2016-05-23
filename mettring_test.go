package mettring

import (
	"fmt"
	"time"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ChristianKniep/QNIBCollect/src/fullerite/metric"
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
	item, ok := mr.Peek(now.Unix())
	assert.True(t, ok, "nothing in the slice")
	assert.Equal(t, m, item[0])
}


func TestMettringValues(t *testing.T) {
	mr := New(2)
	_, ok := mr.Values()
	exp := []metric.Metric{}
	var m metric.Metric
	assert.False(t, ok, "Values() returned non-empty list")
	for i := 0; i < 5; i++ {
	  m = metric.New(fmt.Sprintf("m%d", i))
		exp = append(exp, m)
		mr.Enqueue(m)
		time.Sleep(500 * time.Millisecond)
	}
	slice, ok := mr.Values()
	assert.True(t, ok, "Values() returned empty list")
	assert.Equal(t, exp, slice)
}
