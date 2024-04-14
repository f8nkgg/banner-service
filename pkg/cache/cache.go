package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type Eviction struct {
	Key   compositeKey
	Value interface{}
}
type compositeKey struct {
	part1 int32
	part2 int32
}
type MemoryCache struct {
	values          map[compositeKey]*cacheEntry
	freqs           *list.List
	capacity        int
	sub             int
	len             int
	EvictionChannel chan<- Eviction
	sync.Mutex
}

type cacheEntry struct {
	key      compositeKey
	value    map[string]interface{}
	freqNode *list.Element
	expiry   time.Time
}

type listEntry struct {
	entries map[*cacheEntry]byte
	freq    int
}

func NewMemoryCache(capacity, sub int) *MemoryCache {
	return &MemoryCache{
		values:   make(map[compositeKey]*cacheEntry),
		freqs:    list.New(),
		capacity: capacity,
		sub:      sub,
	}
}

func (c *MemoryCache) Get(key1, key2 int32) (map[string]interface{}, error) {
	c.Lock()
	defer c.Unlock()
	key := compositeKey{part1: key1, part2: key2}
	if e, ok := c.values[key]; ok {
		if e.isExpired() {
			delete(c.values, key)
			c.removeEntry(e.freqNode, e)
			c.len--
			return nil, errors.New("cache expired")
		}
		c.increment(e)
		return e.value, nil
	}
	return nil, errors.New("cache miss")
}
func (e *cacheEntry) isExpired() bool {
	return time.Now().After(e.expiry)
}
func (c *MemoryCache) Set(key1, key2 int32, value map[string]interface{}, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	key := compositeKey{part1: key1, part2: key2}
	if e, exists := c.values[key]; exists {
		e.value = value
		e.expiry = time.Now().Add(ttl)
		c.increment(e)
	} else {
		e := &cacheEntry{
			key:    key,
			value:  value,
			expiry: time.Now().Add(ttl),
		}
		c.values[key] = e
		c.increment(e)
		c.len++
		if c.len > c.capacity {
			c.evict(c.sub)
		}
	}
}

func (c *MemoryCache) evict(count int) int {
	var evicted int
	for i := 0; i < count; {
		if place := c.freqs.Front(); place != nil {
			for entry := range place.Value.(*listEntry).entries {
				if i < count {
					if c.EvictionChannel != nil {
						c.EvictionChannel <- Eviction{
							Key:   entry.key,
							Value: entry.value,
						}
					}
					delete(c.values, entry.key)
					c.removeEntry(place, entry)
					evicted++
					c.len--
					i++
				}
			}
		}
	}
	return evicted
}

func (c *MemoryCache) increment(e *cacheEntry) {
	currentPlace := e.freqNode
	var nextFreq int
	var nextPlace *list.Element
	if currentPlace == nil {
		nextFreq = 1
		nextPlace = c.freqs.Front()
	} else {
		nextFreq = currentPlace.Value.(*listEntry).freq + 1
		nextPlace = currentPlace.Next()
	}

	if nextPlace == nil || nextPlace.Value.(*listEntry).freq != nextFreq {
		li := &listEntry{
			entries: make(map[*cacheEntry]byte),
			freq:    nextFreq,
		}
		if currentPlace != nil {
			nextPlace = c.freqs.InsertAfter(li, currentPlace)
		} else {
			nextPlace = c.freqs.PushFront(li)
		}
	}
	e.freqNode = nextPlace
	nextPlace.Value.(*listEntry).entries[e] = 1
	if currentPlace != nil {
		c.removeEntry(currentPlace, e)
	}
}

func (c *MemoryCache) removeEntry(place *list.Element, entry *cacheEntry) {
	entries := place.Value.(*listEntry).entries
	delete(entries, entry)
	if len(entries) == 0 {
		c.freqs.Remove(place)
	}
}
