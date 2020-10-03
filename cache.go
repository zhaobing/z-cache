package z_cache

import (
	lru "github.com/zhaobing/z-cache/lru"
	"sync"
)

type cache interface {
	add(key string, value ByteView)
	get(key string) (value ByteView, ok bool)
}

type mCache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (m *mCache) add(key string, value ByteView) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.lru == nil { //lazy initialization
		m.lru = lru.New(m.cacheBytes, nil)
	}
	m.lru.Add(key, value)
}

func (m *mCache) get(key string) (value ByteView, ok bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.lru == nil {
		return
	}

	if v, ok := m.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
