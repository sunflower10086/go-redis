package database

import (
	"github.com/sunflower10086/go-redis/interface/database"
	"github.com/sunflower10086/go-redis/interface/resp"
	"github.com/sunflower10086/go-redis/resp/reply"
)

func init() {
	RegisterCmd("Get", execGet, 2)
	RegisterCmd("Set", execSet, 3)
	RegisterCmd("GetSet", execGetSet, 3)
	RegisterCmd("SetNX", execSetNx, 3)
	RegisterCmd("StrLen", execStrLen, 2)
}

// GET
func execGet(db *DB, args CmdLine) resp.Reply {
	key := string(args[0])
	entity, ok := db.GetEntity(key)
	if !ok {
		return reply.NewNullBulkReply()
	}
	bytes := entity.Data.([]byte)
	return reply.NewBulkReply(bytes)
}

// SET
func execSet(db *DB, args CmdLine) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity := database.DataEntity{Data: val}
	db.PutEntity(key, &entity)
	return reply.NewOkReply()
}

// SETNX
func execSetNx(db *DB, args CmdLine) resp.Reply {
	key := string(args[0])
	val := args[1]
	entity := database.DataEntity{Data: val}
	result := db.PutIfAbsent(key, &entity)
	return reply.NewIntReply(int64(result))
}

// GETSET 设置把给定key的值改变成传入的值，并返回旧的值
func execGetSet(db *DB, args CmdLine) resp.Reply {
	key, newVal := string(args[0]), string(args[1])
	entity, ok := db.GetEntity(key)

	db.PutEntity(key, &database.DataEntity{Data: newVal})

	if !ok {
		return reply.NewNullBulkReply()
	}
	oldVal := entity.Data.([]byte)
	return reply.NewBulkReply(oldVal)
}

// STRLEN
func execStrLen(db *DB, args CmdLine) resp.Reply {
	key := string(args[0])
	entity, ok := db.GetEntity(key)
	if !ok {
		return reply.NewNullBulkReply()
	}
	bytes := entity.Data.([]byte)
	return reply.NewIntReply(int64(len(bytes)))
}
