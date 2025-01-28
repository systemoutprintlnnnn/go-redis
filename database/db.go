package database

import (
	"go-redis/datastruct/dict"
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
	"strings"
)

type DB struct {
	index  int
	data   dict.Dict
	addAof func(CmdLine)
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func NewDB() *DB {
	return &DB{
		data: dict.NewSyncDict(),
		addAof: func(line CmdLine) {
			//空方法，为了启动项目时调用loadAof的时候不会写入AOF文件
		},
	}
}

func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command '" + cmdName + "'")
	}
	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	fun := cmd.executor
	return fun(db, cmdLine[1:])
}

// 参数个数固定：SET K V -> arity = 3
// 参数个数不固定，arity表示最小参数数量，这里符号仅为标识：EXISTS k1 k2 -> arity = -2
func validateArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

// GetEntity GET k
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity := raw.(*database.DataEntity)
	return entity, true
}

// PutEntity SET k v
func (db *DB) PutEntity(key string, entity *database.DataEntity) {
	db.data.Put(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

//func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
//	return db.data.PutIfExists(key, entity)
//}

func (db *DB) RemoveEntity(key string) int {
	return db.data.Remove(key)
}

func (db *DB) RemoveEntities(keys ...string) (affected int) {
	affected = 0
	for _, key := range keys {
		affected += db.data.Remove(key)
	}
	return affected
}

func (db *DB) Exist(key string) bool {
	_, ok := db.data.Get(key)
	return ok
}

func (db *DB) Exists(keys ...string) (count int) {
	count = 0
	for _, key := range keys {
		if db.Exist(key) {
			count++
		}
	}
	return count
}

func (db *DB) Flush() {
	db.data.Clear()
}
