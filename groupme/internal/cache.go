package internal

import (
	"sync"
)

// BotIDCache for caching chat to bot relationships
type BotIDCache interface {
	Set(groupID string, bot Bot)
	Get(groupID string) (bot Bot, ok bool)
	Clear()
}

// NewCache returns new BotIDCache
func NewCache() BotIDCache {
	c := &simpleMapCache{
		entries: make(map[string]Bot),
		mutex:   &sync.RWMutex{},
	}

	return c
}

func (c *simpleMapCache) Set(groupID string, bot Bot) {
	c.mutex.Lock()
	c.entries[groupID] = bot
	c.mutex.Unlock()
}

func (c *simpleMapCache) Get(groupID string) (Bot, bool) {
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
	entries map[string]Bot
	mutex   *sync.RWMutex
}
