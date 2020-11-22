package z_cache

// 根据传入的key选择相应节点的PeerGetter
type PeerPicker interface {
	SelectPeer(key string) (peer PeerGetter, ok bool)
}

// 从对应的group中查找缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
