package lru

import (
	"container/list"
)

//lru 缓存淘汰策略
//Cache is a LRU cache, It is not safe for concurrent access
type Cache struct {
	//Maximum memory allowed
	maxBytes int64
	//Memory currently used
	nBytes int64
	//Two-way linked-list save value
	ll *list.List
	//Dictionary, key is string ,value is a element of list
	cache map[string]*list.Element
	//Callback function when records are removed
	OnEvicted func(key string, value Value)
}

//Data type of node in  two-way linked-list.
//The advantage of saving the key corresponding to each value in the linked list is that when eliminating the leader node,can use key to remove the mapping from the dictionary.
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvictedFuc func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvictedFuc,
	}
}

func (c *Cache) Add(key string, value Value) {
	if elem, ok := c.cache[key]; ok { //update
		c.ll.MoveToFront(elem)
		oldEntry := elem.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(oldEntry.value.Len())
		oldEntry.value = value
	} else { //add
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

//RemoveOldest remove the oldest elem from the map and list
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		entry := ele.Value.(*entry)
		delete(c.cache, entry.key)
		c.nBytes -= int64(len(entry.key)) + int64(entry.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(entry.key, entry.value)
		}
	}

}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
