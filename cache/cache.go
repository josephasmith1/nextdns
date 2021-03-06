package cache

import (
	"sync/atomic"

	lru "github.com/hashicorp/golang-lru"
)

type Cache struct {
	*lru.ARCCache

	hit, miss uint64
}

type Stats struct {
	Hit  uint64 `json:"hit"`
	Miss uint64 `json:"miss"`
}

func New(size int) (*Cache, error) {
	arc, err := lru.NewARC(size)
	if err != nil {
		return nil, err
	}
	return &Cache{ARCCache: arc}, nil
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	if c == nil {
		return nil, false
	}
	value, ok = c.ARCCache.Get(key)
	if ok {
		atomic.AddUint64(&c.hit, 1)
	} else {
		atomic.AddUint64(&c.miss, 1)
	}
	return value, ok
}

func (c *Cache) Stats() Stats {
	return Stats{
		Hit:  atomic.LoadUint64(&c.hit),
		Miss: atomic.LoadUint64(&c.miss),
	}
}
