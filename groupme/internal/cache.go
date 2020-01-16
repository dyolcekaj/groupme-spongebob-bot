package internal

import (
	"sync"
)

type BotIdCache interface {
	Set(groupId string, botId string)
	Get(groupId string) (botId string, ok bool)
	Clear()
}

func NewCache() BotIdCache {
	c := &simpleMapCache{
		entries: make(map[string]string),
		mutex:   &sync.RWMutex{},
	}

	return c
}

func (c *simpleMapCache) Set(groupId string, botId string) {
	c.mutex.Lock()
	c.entries[groupId] = botId
	c.mutex.Unlock()
}

func (c *simpleMapCache) Get(groupId string) (string, bool) {
	c.mutex.RLock()
	entry, ok := c.entries[groupId]
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
