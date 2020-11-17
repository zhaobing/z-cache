package consistent_hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFun func(data []byte) uint32

type HashCircle struct {
	hashFun                HashFun
	replicas               int
	logicalNodeHashCodes   []int
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
			logicalNodeKey := []byte(strconv.Itoa(i) + nodeName)
			logicalNodeHashCode := int(m.hashFun(logicalNodeKey))
			m.logicalNodeHashCodes = append(m.logicalNodeHashCodes, logicalNodeHashCode)
			m.logicalPhysicalNodeMap[logicalNodeHashCode] = nodeName
		}
	}
	sort.Ints(m.logicalNodeHashCodes)
}
