package internal

import (
	"sync"
)

// BotIDCache for caching chat to bot relationships
type BotIDCache interface {
	Set(groupID string, botID string)
	Get(groupID string) (botID string, ok bool)
	Clear()
}

// NewCache returns new BotIDCache
func NewCache() BotIDCache {
	c := &simpleMapCache{
		entries: make(map[string]string),
		mutex:   &sync.RWMutex{},
	}

	return c
}

func (c *simpleMapCache) Set(groupID string, botID string) {
	c.mutex.Lock()
	c.entries[groupID] = botID
	c.mutex.Unlock()
}

func (c *simpleMapCache) Get(groupID string) (string, bool) {
	c.mutex.RLock()
	entry, ok := c.entries[groupID]
	c.mutex.RUnlock()

	return entry, ok
}

func (c *simpleMapCache) Clear() {
	c.mutex.Lock()
	for k := range c.entries {
		delete(c.entries, k)
	}
	c.mutex.Unlock()
}

type simpleMapCache struct {
	entries map[string]string
	mutex   *sync.RWMutex
}
