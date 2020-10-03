package z_cache

type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return v.Len()
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
