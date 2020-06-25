package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestCache_RemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := String("value1"), String("value2"), String("v3")
	cap := len(k1+k2) + v1.Len() + v2.Len()
	lru := New(int64(cap), nil)
	lru.Add(k1, v1)
	lru.Add(k2, v2)

	get, ok := lru.Get(k1)
	assert(t, ok, true)
	assert(t, get, v1)

	lru.Add(k3, v3)
	get, ok = lru.Get(k2)
	assert(t, ok, false)
	assert(t, get, nil)

}

func TestCache_Get(t *testing.T) {
	lru := New(int64(0), nil)
	k, v := "key1", String("value1")
	lru.Add(k, v)
	get, ok := lru.Get(k)
	assert(t, ok, true)
	assert(t, get, v)
	get, ok = lru.Get("key2")
	assert(t, ok, false)
	assert(t, get, nil)
}

func TestCache_OnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}

	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}

}

func assert(t *testing.T, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("got is wrong, got %q,but want %q", got, want)
	}
}
