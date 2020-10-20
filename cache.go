package z_cache

import (
	lru2 "github.com/zhaobing/z-cache/lru"
	"sync"
)

type cache interface {
	add(key string, value ByteView)
	get(key string) (value ByteView, ok bool)
}

type mCache struct {
	mu            sync.Mutex
	lru           *lru2.Cache
	maxLimitBytes int64
}

func (m *mCache) add(key string, value ByteView) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.lru == nil {
		m.lru = lru2.New(m.maxLimitBytes, nil)
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
