package z_cache

import (
	lru2 "github.com/zhaobing/z-cache/lru"
	"sync"
)

//并发控制封装
type cache interface {
	add(key string, value ByteView)
	get(key string) (value ByteView, ok bool)
}

type mCache struct {
	mutex         sync.Mutex
	lru           *lru2.Cache
	maxLimitBytes int64
}

func (m *mCache) add(key string, value ByteView) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//延迟初始化,是否可以使用sync.one优化?
	if m.lru == nil {
		m.lru = lru2.New(m.maxLimitBytes, nil)
	}
	m.lru.Add(key, value)
}

func (m *mCache) get(key string) (value ByteView, ok bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.lru == nil {
		return
	}

	if v, ok := m.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
