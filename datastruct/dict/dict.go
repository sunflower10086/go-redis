package dict

type Consumer func(key string, val any) bool

type Dict interface {
	Get(key string) (val any, exists bool)
	Len() int
	Put(key string, val any) (result int)
	// PutIfAbsent SETNX, 如果没有再set进去
	PutIfAbsent(key string, val any) (result int)
	// PutIfExists 有这个key的时候再往里set，没有就不set
	PutIfExists(key string, val any) (result int)
	// Remove 删除一个key
	Remove(key string) (result int)
	// ForEach 遍历所有的kv
	ForEach(consumer Consumer)
	// Keys 列出所有的key
	Keys() []string
	// RandomKeys 随机返回limit个key
	RandomKeys(limit int) []string
	// RandomDistinckKeys 随机返回limit个不同的key
	RandomDistinckKeys(limit int) []string
	// Clear 清空数据
	Clear()
}
