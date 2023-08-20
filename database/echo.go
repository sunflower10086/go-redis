package database

import (
	"github.com/sunflower10086/go-redis/interface/resp"
	"github.com/sunflower10086/go-redis/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

func (e *EchoDatabase) Exec(client resp.Connection, data [][]byte) resp.Reply {

	return reply.NewMultiBulkReply(data)
}

func (e *EchoDatabase) Close() {

}

func (e *EchoDatabase) AfterClientClose(client resp.Connection) {

}
