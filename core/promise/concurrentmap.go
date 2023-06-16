package promise

import (
	"sync"
)

type ConcurrencyMap struct {
	m  sync.Map
	mu sync.Mutex
}

func NewConcurrencyMap() *ConcurrencyMap {
	return &ConcurrencyMap{
		m: sync.Map{},
	}
}

func (c *ConcurrencyMap) Put(key string, val interface{}) {
	// c.mu.Lock()
	// defer c.mu.Unlock()
	c.m.Store(key, val)
}

func (c *ConcurrencyMap) GetAndClear() map[string]interface{} {

	// c.mu.Lock()
	// defer c.mu.Unlock()

	result := make(map[string]interface{}, 0)

	c.m.Range(func(key interface{}, value interface{}) bool {
		result[key.(string)] = value.([]interface{})
		c.m.Delete(key)
		return true
	})

	return result

}

func (c *ConcurrencyMap) Get(key string) interface{} {
	val, ok := c.m.Load(key)
	if ok {
		return val.(interface{})
	}
	return nil
}

func (c *ConcurrencyMap) ToMap() map[string]interface{} {
	result := make(map[string]interface{}, 0)
	c.m.Range(func(key interface{}, value interface{}) bool {
		result[key.(string)] = value
		return true
	})
	return result
}
