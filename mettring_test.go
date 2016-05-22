package mettring

import (
	//"fmt"
	"time"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestInt struct {
	key string
	val string
	time time.Time
}

func (t TestInt) Time() time.Time {
	return t.time
}

// SetTime overwrites the time
func (t TestInt) SetTime(v time.Time) {
	t.time = v
}

// New returns a new ring
func NewTest(key string, val string) TestInt {
	return TestInt{
		key: key,
		val: val,
		time: time.Now(),
	}
}

func TestNewMettring(t *testing.T) {
	mr := New(300)
	assert.Equal(t, mr.retention, 300)
}

func TestMettringEnque(t *testing.T) {
	mr := New(300)
	ti := NewTest("key", "val0")
	now := time.Now()
	ti.SetTime(now)
	mr.Enqueue(ti)
	//item, ok := mr.Peek(now.Unix())
	//assert.True(t, ok, "nothing in the slice")
	//assert.Equal(t, ti, item)
}
