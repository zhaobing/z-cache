package z_cache

import (
	"github.com/zhaobing/bingo/utils/test"
	"testing"
)

func Test_mCache(t *testing.T) {

	key := "leb"
	value := NewByteViewByString("ronNG")
	mCache := &mCache{maxLimitBytes: 66666}
	get1, _ := mCache.get(key)
	if len(get1.b) != 0 {
		t.Errorf("get1 should empty,but len is %d", len(get1.b))
	}
	mCache.add(key, value)
	get2, _ := mCache.get(key)
	test.AssertStrings(t, get2.String(), value.String())
}
