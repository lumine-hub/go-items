package geecache

/*
*
ByteView 缓存值的抽象与封装(类似于序列化)
*/
type ByteView struct {
	data []byte
}

func (b ByteView) Len() int {
	return len(b.data)
}

// ByteSlice 获取副本
func (b ByteView) ByteSlice() []byte {
	c := make([]byte, len(b.data))
	copy(b.data, c)
	return c
}

func (b ByteView) String() string {
	return string(b.data)
}
