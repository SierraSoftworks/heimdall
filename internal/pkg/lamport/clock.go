package lamport

import (
	"sync"
)

type Clock struct {
	time uint64

	m sync.RWMutex
}

func NewClock(initial uint64) *Clock {
	return &Clock{
		time: initial,
	}
}

func (c *Clock) Next() uint64 {
	c.m.Lock()
	defer c.m.Unlock()

	c.time = c.time + 1
	return c.time
}

func (c *Clock) Update(time uint64) uint64 {
	c.m.Lock()
	defer c.m.Unlock()

	if time > c.time {
		c.time = time + 1
	} else {
		c.time = c.time + 1
	}

	return c.time
}

func (c *Clock) Current() uint64 {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.time
}

func (c *Clock) NotAfter(time uint64) bool {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.time < time
}
