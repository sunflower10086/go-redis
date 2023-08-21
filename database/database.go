package database

import (
	"strconv"
	"strings"

	"github.com/sunflower10086/go-redis/config"
	"github.com/sunflower10086/go-redis/interface/resp"
	"github.com/sunflower10086/go-redis/lib/logger"
	"github.com/sunflower10086/go-redis/resp/reply"
)

type Database struct {
	dbSet []*DB
}

func NewDatabase() *Database {
	var database Database
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := 0; i < config.Properties.Databases; i++ {
		db := NewDB()
		db.index = i
		database.dbSet[i] = db
	}

	return &database
}

// args
// set k v
// get k
// select 2 只有select是在这一层做的
func (d *Database) Exec(conn resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	// 选择数据库
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.NewArgNumErrReply("select")
		}
		return execSelectDB(conn, d, args[1:])
	}

	index := conn.GetDBIndex()
	nowDB := d.dbSet[index]
	return nowDB.Exec(conn, args)
}

func (d *Database) Close() {

}

func (d *Database) AfterClientClose(connection resp.Connection) {

}

func execSelectDB(conn resp.Connection, database *Database, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.NewErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(database.dbSet) {
		return reply.NewErrReply("ERR DB index is out if range")
	}

	conn.SelectDB(dbIndex)
	return reply.NewOkReply()
}
