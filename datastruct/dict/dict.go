package dict

type Consumer func(key string, val interface{}) bool

// 实现Redis的Dict数据结构
type Dict interface {
	// 得到key对应的值
	Get(key string) (interface{}, bool)
	// 返回字典中的元素个数
	Len() int
	//存入键值对
	//Return: 影响行数
	Put(key string, val interface{}) (result int)
	PutIfAbsent(key string, val interface{}) (result int)
	//PutIfExists(key string, val interface{}) (result int)
	Remove(key string) (result int)
	ForEach(consumer Consumer)
	//列出所有key
	Keys() []string
	//列出随机n个key
	RandomKeys(n int) []string
	//返回n个不重复的key
	RandomDistinctKeys(n int) []string
	Clear()
}
