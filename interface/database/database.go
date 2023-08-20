package database

import "github.com/sunflower10086/go-redis/interface/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(resp.Connection, [][]byte) resp.Reply
	Close()
	AfterClientClose(connection resp.Connection)
}

type DataEntity struct {
	Data interface{}
}
