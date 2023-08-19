package resp

// Reply 服务端对客户端的回复
type Reply interface {
	ToBytes() []byte
}
