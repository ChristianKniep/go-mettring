package mettring

import (
	//"fmt"
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
