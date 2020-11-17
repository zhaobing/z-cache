package z_cache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func key2Bytes(key string) ([]byte, error) {
	return []byte(key), nil
}

func TestGetter1(t *testing.T) {
	t.Helper()
	t.Run("外部命名函数", func(t *testing.T) {
		//cast function to  interface
		//借助 GetterFunc 的类型转换，将一个匿名回调函数转换成了接口 f Getter。
		//调用该接口的方法 f.Get(key string)，实际上就是在调用匿名回调函数
		//定义一个函数类型 F，并且实现接口 A 的方法，然后在这个方法中调用自己。这是 Go 语言中将其他函数（参数返回值定义与 F 一致）转换为接口 A 的常用技巧。
		var f Getter = GetterFunc(key2Bytes)

		expect := []byte("key")
		get, err := f.Get("key")

		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(get, expect) {
			t.Errorf("callback failed")
		}
	})

	t.Run("内部函数赋值", func(t *testing.T) {
		//cast function to  interface

		var key2BytesFunc = func(key string) ([]byte, error) {
			return []byte(key), nil
		}

		var f Getter = GetterFunc(key2BytesFunc)

		expect := []byte("key")
		get, err := f.Get("key")

		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(get, expect) {
			t.Errorf("callback failed")
		}
	})

	t.Run("匿名函数", func(t *testing.T) {

		key := "hello"
		expect := []byte(key)

		get, err := GetterFunc(func(key string) ([]byte, error) {
			return []byte(key), nil
		}).Get(key)

		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(get, expect) {
			t.Errorf("callback error")
		}

	})

}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "630",
	"Sam":  "630",
}

func TestGet(t *testing.T) {
	t.Helper()

	loadCounts := make(map[string]int, len(db))

	scoresGroup := NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s is not exist", key)
	}))

	t.Run("zCacheGet-连续执行2次获取相同key", func(t *testing.T) {

		for k, v := range db {
			if view, err := scoresGroup.Get(k); err != nil || view.String() != v {
				t.Fatal("fail to get value from cache")
			}

			if _, err := scoresGroup.Get(k); err != nil || loadCounts[k] > 1 {
				t.Fatalf("get data from cache repeate error ,key is %s", k)
			}
		}

	})

	t.Run("zCacheGet-获取不存在的key", func(t *testing.T) {
		if view, err := scoresGroup.Get("unknown"); err == nil {
			t.Fatalf("the value of unknown should be empty,but %s got", view)
		}
	})

}
