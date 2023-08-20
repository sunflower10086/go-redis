package resp

// Connection 对应redis的每一个连接
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}
