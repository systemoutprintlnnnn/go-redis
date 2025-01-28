package dict

import "sync"

type SyncDict struct {
	m sync.Map
}

func NewSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (interface{}, bool) {
	val, ok := dict.m.Load(key)
	return val, ok
}

func (dict *SyncDict) Len() int {
	length := 0
	dict.m.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	value, ok := dict.m.Load(key)
	dict.m.Store(key, val)
	if ok && val == value {
		return 0
	}
	return 1
}

func (dict *SyncDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, ok := dict.m.Load(key)
	if !ok {
		dict.m.Store(key, val)
		return 1
	}
	return 0
}

//func (dict *SyncDict) PutIfExists(key string, val interface{}) (result int) {
//	_, ok := dict.m.Load(key)
//	if ok {
//		dict.m.Store(key, val)
//		return 1
//	}
//	return 0
//}

func (dict *SyncDict) Remove(key string) (result int) {
	_, ok := dict.m.Load(key)
	dict.m.Delete(key)
	if ok {
		return 1
	}
	return 0
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value interface{}) bool {
		consumer(key.(string), value)
		return true
	})
}

func (dict *SyncDict) Keys() []string {
	result := make([]string, dict.Len())
	i := 0
	dict.m.Range(func(key, value interface{}) bool {
		result[i] = key.(string)
		i++
		return true
	})
	return result
}

// 返回随机n个key（可重复）
func (dict *SyncDict) RandomKeys(n int) []string {
	result := make([]string, n)
	for i := 0; i < n; i++ {
		dict.m.Range(func(key, value interface{}) bool {
			result[i] = key.(string)
			return false
		})
	}
	return result
}

func (dict *SyncDict) RandomDistinctKeys(n int) []string {
	result := make([]string, n)
	i := 0
	dict.m.Range(func(key, value interface{}) bool {
		result[i] = key.(string)
		i++
		return i < n
	})
	return result
}

func (dict *SyncDict) Clear() {
	*dict = *NewSyncDict()
}
