package cache

import "sync"

type Etag struct {
	// Key stores etag value
	Key string
	// Data is HTTP response
	Data []byte
}

// Storage interface for HTTP response caching
type Storage interface {
	Get(key string) *Etag
	Set(key string, etag Etag)
	Delete(key string)
}

// InMemoryStore implements in-memory cache
type InMemoryStore struct {
	cache map[string]*Etag
	sync.RWMutex
}

func (im *InMemoryStore) Get(key string) *Etag {
	im.RLock()
	defer im.RUnlock()
	if val, ok := im.cache[key]; ok {
		return val
	}
	return nil
}

func (im *InMemoryStore) Set(key string, etag Etag) {
	im.Lock()
	defer im.Unlock()
	im.cache[key] = &etag
}

func (im *InMemoryStore) Delete(key string) {
	im.Lock()
	defer im.Unlock()
	delete(im.cache, key)
}

func NewStorage() Storage {
	return &InMemoryStore{
		cache:   make(map[string]*Etag),
		RWMutex: sync.RWMutex{},
	}
}
