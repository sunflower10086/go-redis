package database

import (
	"github.com/sunflower10086/go-redis/interface/resp"
	"github.com/sunflower10086/go-redis/resp/reply"
)

func init() {
	RegisterCmd("ping", Ping, 1)
}

func Ping(db *DB, args CmdLine) resp.Reply {
	return reply.NewPongReply()
}
