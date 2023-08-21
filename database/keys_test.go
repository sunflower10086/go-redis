package database

import (
	"fmt"
	"testing"

	"github.com/sunflower10086/go-redis/interface/database"
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

func TestKeys(t *testing.T) {
	testDB := NewDB()

	exec := testDB.PutEntity("key", &database.DataEntity{Data: [][]byte{[]byte("val")}})
	fmt.Println(exec)
	entity, _ := testDB.GetEntity("key")
	fmt.Println(string(entity.Data.([][]byte)[0]))

	keys := execKeys(testDB, [][]byte{[]byte("*")})
	fmt.Println(string(keys.ToBytes()))
}
