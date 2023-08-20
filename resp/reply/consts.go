package reply

// PongRelpy PONG的结构体
type PongRelpy struct {
}

var pongbytes = []byte("+PONG\r\n")

func (p *PongRelpy) ToBytes() []byte {
	return pongbytes
}

func NewPongReply() *PongRelpy {
	return &PongRelpy{}
}

// OkReply 返回OK
type OkReply struct{}

var okBytes = []byte("+OK\r\n")

func (o *OkReply) ToBytes() []byte {
	return okBytes
}

func NewOkReply() *OkReply {
	return &OkReply{}
}

// NullBulkReply 空回复
type NullBulkReply struct{}

var nullBulkBytes = []byte("$-1\r\n")

func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func NewNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

// EmptyMultiBulkReply 一个空数组
type EmptyMultiBulkReply struct{}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

func NewEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

// NoReply 真空，什么都没有
type NoReply struct{}

var noBytes = []byte("")

func (n *NoReply) ToBytes() []byte {
	return noBytes
}

func NewNoReply() *NoReply {
	return &NoReply{}
}
