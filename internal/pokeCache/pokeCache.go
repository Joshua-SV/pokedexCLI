package pokeCache

import (
	"fmt"
	"sync"
	"time"
)

type cacheType interface {
	Add(key string, val []byte)
	Get(key string) ([]byte, bool)
	reapLoop()
}

// Cache provides thread-safe access to a timed cache
type Cache struct {
	table    map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		table:    make(map[string]cacheEntry, 0),
		interval: interval,
	}
	// create a go routine to update table entries when expired
	go cache.reapLoop()

	return cache
}

// Add adds a new entry to the cache
func (c *Cache) Add(key string, val []byte) {
	// lock for concurrency safety
	c.mu.Lock()
	defer c.mu.Unlock()

	// add the data into the cache
	c.table[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

// Get retrieves an entry from the cache if it's not expired
func (c *Cache) Get(key string) ([]byte, bool) {
	// read lock incase of any writing go routines
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, okay := c.table[key]
	if okay == false {
		fmt.Println("Failed Cache Hit!!!")
		return nil, false
	}

	// delete data if it passed interval time
	if time.Since(entry.createdAt) > c.interval {
		delete(c.table, key)
		return nil, false
	}

	// Update access time (sliding expiration) Least Resently Used (LRU) policy
	entry.createdAt = time.Now()
	c.table[key] = entry

	fmt.Println("Successful Cache Hit!!!")
	return entry.val, true
}

// reapLoop removes expired entries periodically
func (c *Cache) reapLoop() {
	// create a counter which periodically ticks based on the interval used
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	// .c is a channel by which the ticks are delivered
	for range ticker.C {
		c.mu.Lock()
		currentTime := time.Now()

		// loop through the cache table to check which data expired
		for key, entry := range c.table {
			// if data is older than the interval allowed, delete it
			if currentTime.Sub(entry.createdAt) > c.interval {
				delete(c.table, key)
			}
		}
		c.mu.Unlock()
	}
}
