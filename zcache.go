package z_cache

import (
	"fmt"
	"log"
	"sync"
)

//负责与外部交互，控制缓存存储和获取的主流程

// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache mCache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, maxLimitBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: mCache{
			maxLimitBytes: maxLimitBytes,
		},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[zcache] hit")
		return v, nil
	}

	//key没有对应的缓存值，需要加载
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneByteView(bytes)}
	g.populateCache(key, value)
	return value, nil

}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
