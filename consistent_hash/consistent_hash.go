package consistent_hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFun func(data []byte) uint32

type HashCircle struct {
	//哈希函数
	hashFun HashFun
	//虚拟节点数量
	replicas int
	//逻辑节点的hashCode
	logicalNodeHashCodes []int
	//逻辑节点hashCode与物理节点名称的映射
	logicalPhysicalNodeMap map[int]string
}

func New(replicas int, fn HashFun) *HashCircle {
	m := &HashCircle{
		hashFun:                fn,
		replicas:               replicas,
		logicalPhysicalNodeMap: make(map[int]string),
	}

	if m.hashFun == nil {
		m.hashFun = crc32.ChecksumIEEE
	}
	return m
}

//AddPhysicalNode  add physical node to the hash circle
func (m *HashCircle) AddPhysicalNode(nodeNames ...string) {
	for _, nodeName := range nodeNames {
		for i := 0; i < m.replicas; i++ {
			//logicalNodeKey := []byte(strconv.Itoa(i) + "-" + nodeName)
			logicalNodeKey := []byte(strconv.Itoa(i) + nodeName)
			logicalNodeHashCode := int(m.hashFun(logicalNodeKey))
			m.logicalNodeHashCodes = append(m.logicalNodeHashCodes, logicalNodeHashCode)
			m.logicalPhysicalNodeMap[logicalNodeHashCode] = nodeName
		}
	}
	sort.Ints(m.logicalNodeHashCodes)
}

//GetPhysicalNode 根据Key，获取哈希环上的最近的节点
func (m *HashCircle) GetPhysicalNode(key string) string {
	if len(m.logicalNodeHashCodes) == 0 {
		return ""
	}

	hashCode := int(m.hashFun([]byte(key)))

	idx := sort.Search(len(m.logicalNodeHashCodes), func(i int) bool {
		return m.logicalNodeHashCodes[i] >= hashCode
	})

	idx = idx % len(m.logicalNodeHashCodes)
	logicalNodeHash := m.logicalNodeHashCodes[idx]
	return m.logicalPhysicalNodeMap[logicalNodeHash]
}
