package cache

import (
	"sync"
	"time"
)

const shardSize = 256

// Cache struct
type Cache struct {
	buckets [shardSize]*bucket
}

// New Cache with size
func New(size int) *Cache {
	c := new(Cache)
	bucketSize := size / shardSize
	for i := 0; i < shardSize; i++ {
		c.buckets[i] = newBucket(bucketSize)
	}
	return c
}

func newBucket(size int) *bucket {
	return &bucket{
		droplets: make(map[uint64]*droplet),
		size:     size,
	}
}

// Set cache
func (c *Cache) Set(key uint64, v interface{}, ttl uint32) {
	c.buckets[key&(shardSize-1)].set(key, v, ttl)
}

// Get cache
func (c *Cache) Get(key uint64) (interface{}, bool) {
	return c.buckets[key&(shardSize-1)].get(key)
}

// Del cache
func (c *Cache) Del(key uint64) {
	c.buckets[key&(shardSize-1)].del(key)
}

// Len return the length of current caches
func (c *Cache) Len() int {
	l := 0
	for _, b := range c.buckets {
		l += b.len()
	}
	return l
}

type bucket struct {
	droplets map[uint64]*droplet
	size     int

	sync.RWMutex
}

func (b *bucket) set(k uint64, v interface{}, ttl uint32) {
	if len(b.droplets)+1 > b.size {
		b.evict()
	}

	b.Lock()
	b.droplets[k] = &droplet{
		value:     v,
		expiredAt: uint32(time.Now().Unix()) + ttl,
	}
	b.Unlock()
}

func (b *bucket) get(k uint64) (interface{}, bool) { // lazy ttl
	b.RLock()
	defer b.RUnlock()
	d, ok := b.droplets[k]
	if !ok {
		return nil, false
	}
	if d.isExpired() {
		b.RUnlock()
		b.del(k)
		b.RLock()
		return nil, false
	}
	return d.value, ok
}

func (b *bucket) del(k uint64) {
	b.Lock()
	delete(b.droplets, k)
	b.Unlock()
}

func (b *bucket) len() int {
	b.RLock()
	l := len(b.droplets)
	b.RUnlock()
	return l
}

func (b *bucket) evict() {
	b.Lock()
	defer b.Unlock()
	key := -1
	for k := range b.droplets {
		key = int(k)
		break
	}
	if key == -1 {
		// empty cache
		return
	}
	delete(b.droplets, uint64(key))
}

type droplet struct {
	value     interface{}
	expiredAt uint32
}

func (d *droplet) isExpired() bool {
	return d.expiredAt < uint32(time.Now().Unix())
}
