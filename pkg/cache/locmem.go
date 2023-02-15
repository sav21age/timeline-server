package cache

import (
	"sync"
	"time"
)

type item struct {
	value     interface{}
	createdAt int64
	ttl       int64
}

type LocMem struct {
	cache map[interface{}]*item
	rwm *sync.RWMutex
}

// Cache uses map to store key:value data in-memory.
func NewLocMem() *LocMem {
	c := &LocMem{cache: make(map[interface{}]*item), rwm: &sync.RWMutex{}}
	go c.setTtlTimer()

	return c
}

//--

func (c *LocMem) setTtlTimer() {
	for {
		c.rwm.Lock()
		for k, i := range c.cache {
			if time.Now().Unix()-i.createdAt > i.ttl {
				delete(c.cache, k)
			}
		}
		c.rwm.Unlock()

		<-time.After(time.Second)
	}
}

//--

func (c *LocMem) Set(key, value interface{}, ttl int64) error {
	c.rwm.Lock()
	c.cache[key] = &item{
		value:     value,
		createdAt: time.Now().Unix(),
		ttl:       ttl,
	}
	c.rwm.Unlock()

	return nil
}

//--

func (c *LocMem) Get(key interface{}) (interface{}, error) {
	c.rwm.RLock()
	item, ok := c.cache[key]
	c.rwm.RUnlock()

	if !ok {
		return nil, ErrItemNotFound
	}

	return item.value, nil
}
