package z_cache

import (
	"github.com/zhaobing/bingo/utils/test"
	"testing"
)

func Test_mCache(t *testing.T) {

	key := "leb"
	value := "rong"
	m := &mCache{maxLimitBytes: 66666}
	get1, _ := m.get(key)
	if len(get1.b) != 0 {
		t.Errorf("get1 should empty,but len is %d", len(get1.b))
	}

	m.add(key, NewByteViewByString(value))
	get2, _ := m.get(key)
	test.AssertStrings(t, get2.ToStr(), value)
}
