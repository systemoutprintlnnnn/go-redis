package consistenthash

import (
	"hash/crc32"
	"sort"
)

type HashFunc func(data []byte) uint32
type NodeMap struct {
	hashFunc    HashFunc
	nodeHashes  []int
	nodeHashMap map[int]string
}

func NewNodeMap(hashFunc HashFunc) *NodeMap {
	if hashFunc == nil {
		hashFunc = crc32.ChecksumIEEE
	}
	m := &NodeMap{
		hashFunc:    hashFunc,
		nodeHashMap: make(map[int]string),
	}
	return m
}

func (m *NodeMap) IsEmpty() bool {
	return len(m.nodeHashes) == 0
}

func (m *NodeMap) AddNode(keys ...string) {
	//计算哈希值
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := int(m.hashFunc([]byte(key)))
		m.nodeHashes = append(m.nodeHashes, hash)
		m.nodeHashMap[hash] = key
	}
	sort.Ints(m.nodeHashes)

}

// 根据key决定处理节点
func (m *NodeMap) GetNode(key string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hashFunc([]byte(key)))
	index := sort.Search(len(m.nodeHashes), func(i int) bool {
		return m.nodeHashes[i] >= hash
	})
	if index == len(m.nodeHashes) {
		index = 0
	}
	return m.nodeHashMap[m.nodeHashes[index]]
}
