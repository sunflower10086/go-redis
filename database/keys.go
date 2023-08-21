package database

import (
	"github.com/sunflower10086/go-redis/interface/resp"
	wildcard "github.com/sunflower10086/go-redis/lib/wildacrd"
	"github.com/sunflower10086/go-redis/resp/reply"
)

func init() {
	RegisterCmd("del", execDel, -2)
	RegisterCmd("exists", execExists, -2)
	RegisterCmd("keys", execKeys, 2)
	RegisterCmd("flushdb", execFlushDB, 1)
	RegisterCmd("type", execType, 2)
	RegisterCmd("rename", execRename, 3)
	RegisterCmd("rename", execRenameNx, 3)
}

// DEL k1 k2 k3...
func execDel(db *DB, args CmdLine) resp.Reply {
	keys := make([]string, 0, len(args))
	for _, arg := range args {
		keys = append(keys, string(arg))
	}
	deleted := db.Removes(keys...)
	return reply.NewIntReply(int64(deleted))
}

// EXISTS k1 k2 k3 ... 有几个存在
func execExists(db *DB, args CmdLine) resp.Reply {
	var result int64
	for _, arg := range args {
		_, ok := db.GetEntity(string(arg))
		if ok {
			result++
		}
	}
	return reply.NewIntReply(result)
}

// KEYS 通配符...
func execKeys(db *DB, args CmdLine) resp.Reply {
	pattern, err := wildcard.CompilePattern(string(args[0]))
	if err != nil {
		return reply.NewErrReply("param err")
	}
	result := make([][]byte, 0, 10)
	db.data.ForEach(func(key string, val any) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})

	return reply.NewMultiBulkReply(result)
}

// FLUSHDB
func execFlushDB(db *DB, args CmdLine) resp.Reply {
	db.Flush()
	return reply.NewOkReply()
}

// TYPE key 获得一个key的类型
func execType(db *DB, args [][]byte) resp.Reply {
	key := args[0]
	entity, ok := db.GetEntity(string(key))
	if !ok {
		return reply.NewStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.NewStatusReply("string")
	}
	return reply.NewErrReply("unknown type")
}

// RENAME k1 k2
func execRename(db *DB, args CmdLine) resp.Reply {
	k1, k2 := string(args[0]), string(args[1])
	entity, ok := db.GetEntity(k1)
	if !ok {
		return reply.NewErrReply("no such key")
	}
	db.PutEntity(k2, entity)
	db.Remove(k1)
	return reply.NewOkReply()
}

// RENAMENX 检查k2是否存在
func execRenameNx(db *DB, args CmdLine) resp.Reply {
	k1, k2 := string(args[0]), string(args[1])
	// 检查k2是否存在
	_, ok2 := db.GetEntity(k2)
	if ok2 {
		return reply.NewIntReply(0)
	}

	entity1, ok1 := db.GetEntity(k1)
	if !ok1 {
		return reply.NewErrReply("no such key")
	}
	db.PutEntity(k2, entity1)
	db.Remove(k1)
	return reply.NewIntReply(1)
}
