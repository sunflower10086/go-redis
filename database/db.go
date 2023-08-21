package database

import (
	"strings"

	"github.com/sunflower10086/go-redis/datastruct/dict"
	"github.com/sunflower10086/go-redis/interface/database"
	"github.com/sunflower10086/go-redis/interface/resp"
	"github.com/sunflower10086/go-redis/resp/reply"
)

type DB struct {
	index int
	data  dict.Dict
}

func (db *DB) Close() {

}

func (db *DB) AfterClientClose(connection resp.Connection) {

}

// ExecFunc 执行函数
type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func NewDB() *DB {
	return &DB{
		data: dict.NewSyncDict(),
	}
}

func (db *DB) Exec(conn resp.Connection, cmdLine CmdLine) resp.Reply {
	// PING SET SETNX
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.NewErrReply("ERR unknown command " + cmdName)
	}
	// SET k v
	if !validateArity(cmd.arity, cmdLine) {
		return reply.NewArgNumErrReply(cmdName)
	}
	fun := cmd.exector
	//
	return fun(db, cmdLine[1:])
}

// SET k v -> arity = 3
// EXISTS k1 k2 k3 ... -> arity = -2 (负的最小值，表示可以超过这个值,负号表示边长参数)
func validateArity(arity int, cmdArgs [][]byte) bool {
	argsNum := len(cmdArgs)
	// 判断是不是定长
	if arity >= 0 {
		return argsNum == arity
	} else {
		return argsNum >= -arity
	}
}

func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	val, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}

	entity, _ := val.(*database.DataEntity)

	return entity, true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

func (db *DB) Remove(key string) int {
	return db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) int {
	delSize := 0
	for _, key := range keys {
		delSize += db.data.Remove(key)
	}

	return delSize
}

func (db *DB) Flush() {
	db.data.Clear()
}
