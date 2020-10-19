package z_cache

import (
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
}
