package geecache

import (
	"GeeCache/lru"
	"sync"
)

/*
Cache 封装原来的lru缓存和序列化，增加锁允许并法
value的值从原来的value变成了ByteView
*/
type Cache struct {
	lru        *lru.Cache
	mu         sync.Mutex
	cacheBytes int64
}

func (c *Cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *Cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	res, ok := c.lru.Get(key)
	if ok {
		return res.(ByteView), ok
	}
	return
}
