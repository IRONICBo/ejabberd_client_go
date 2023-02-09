package utils

import (
	"sync/atomic"
)

// 一般不会直接把变量暴露出去，而是通过方法来访问
type counter struct {
	value int64
}

func NewCounter() *counter {
	return &counter{}
}

func (c *counter) Inc() {
	atomic.AddInt64(&c.value, 1)
}

func (c counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}
