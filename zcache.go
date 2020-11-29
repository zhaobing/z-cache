package z_cache

import (
	"fmt"
	"github.com/zhaobing/z-cache/singleflight"
	pb "github.com/zhaobing/z-cache/zcachepb"
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

// Get implements Getter interface function,接口型函数，保证灵活性
func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name       string
	getter     Getter
	mainCache  mCache
	peerPicker PeerPicker
	loader     *singleflight.Group
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
	group := &Group{
		name:   name,
		getter: getter,
		mainCache: mCache{
			maxLimitBytes: maxLimitBytes,
		},
		loader: &singleflight.Group{},
	}
	groups[name] = group
	return group
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

//RegisterPeers 将实现了PeerPicker的HTTPPool注入到Group中
func (g *Group) RegisterPeers(peerPicker PeerPicker) {
	if g.peerPicker != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peerPicker = peerPicker
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

func (g *Group) load(key string) (value ByteView, err error) {
	// each key is only fetched once (either locally or remotely)
	// regardless of the number of concurrent callers.
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peerPicker != nil {
			if peer, ok := g.peerPicker.SelectPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[zcache] Failed to get from peer", err)
			}
		}

		log.Println("g.peerPicker == nil", g.peerPicker == nil)
		return g.getLocally(key)
	})

	if err == nil {
		return view.(ByteView), nil
	}
	return
}

//getFromPeer  从远程节点获取缓存值
func (g *Group) getFromPeer(peerGetter PeerGetter, key string) (ByteView, error) {

	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}

	err := peerGetter.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{
		b: res.Value,
	}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	log.Println("getLocally")
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := CloneByteView(bytes)
	g.populateCache(key, value)
	return value, nil

}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
