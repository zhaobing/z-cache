package z_cache

// 缓存值的抽象与封装
//字节视图，构造一个不变的数据结构，方便缓存复制操作
type ByteView struct {
	b []byte
}

//字节视图的大小,实现了Value接口
func (bv ByteView) Len() int {
	return len(bv.b)
}

//字节视图的拷贝
func (bv ByteView) ByteSlice() []byte {
	return cloneByteView(bv.b)
}

func CloneByteView(b []byte) ByteView {
	bytes := make([]byte, len(b))
	copy(bytes, b)
	return ByteView{b: bytes}
}

func cloneByteView(b []byte) []byte {
	bytes := make([]byte, len(b))
	copy(bytes, b)
	return bytes
}

func (bv ByteView) String() string {
	return string(bv.b)
}

func NewByteViewByString(str string) ByteView {
	bytes := []byte(str)
	return ByteView{b: bytes}
}
