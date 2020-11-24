package z_cache

//PeerPicker 节点选择，根据传入的key选择相应节点的PeerGetter
type PeerPicker interface {
	SelectPeer(key string) (peer PeerGetter, ok bool)
}

//PeerGetter 缓存获取Http客户端，从对应的group中查找key对应缓存值
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
