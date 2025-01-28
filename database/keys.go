package database

import (
	"go-redis/interface/resp"
	"go-redis/lib/utils"
	"go-redis/lib/wildcard"
	"go-redis/resp/reply"
)

// DEL k1 k2...
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	affected := db.RemoveEntities(keys...)
	if affected > 0 {
		db.addAof(utils.ToCmdLine2("DEL", args...))
	}
	return reply.MakeIntReply(int64(affected))
}

// EXISTS k1 k2...
func execExists(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	affected := db.Exists(keys...)
	return reply.MakeIntReply(int64(affected))
}

// FLUSHDB Removes all keys from the current database
// 可能有bug，flush后面接参数
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	db.addAof(utils.ToCmdLine("FLUSHDB"))
	return reply.MakeOkReply()
}

// TYPE k1
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeBulkReply([]byte("none"))
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	}
	//还没有实现的类型
	return reply.MakeUnknownErrReply()
}

// RENAME k1 k2  会覆盖写，将 k1 改名为 k2, v不变
func execRename(db *DB, args [][]byte) resp.Reply {
	key1 := string(args[0])
	key2 := string(args[1])
	entity, ok := db.GetEntity(key1)
	if !ok {
		return reply.MakeErrReply("no such key")
	}
	db.PutEntity(key2, entity)
	db.RemoveEntity(key1)
	return reply.MakeOkReply()
}

// RENAMEX k1 k2  不会覆盖写，将 k1 改名为 k2, v不变
// Returns: affected elements
func execRenameX(db *DB, args [][]byte) resp.Reply {
	key1 := string(args[0])
	key2 := string(args[1])
	entity, ok := db.GetEntity(key1)
	if !ok {
		return reply.MakeErrReply("no such key")
	}
	if db.Exist(key2) {
		return reply.MakeIntReply(0)
	}
	db.PutEntity(key2, entity)
	db.RemoveEntity(key1)
	return reply.MakeIntReply(1)
}

// KEYS
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(result)
}

func init() {
	RegisterCommand("DEL", execDel, -2)
	RegisterCommand("EXISTS", execExists, -2)
	RegisterCommand("FLUSHDB", execFlushDB, 1)
	RegisterCommand("TYPE", execType, 2)
	RegisterCommand("RENAME", execRename, 3)
	RegisterCommand("RENAMENX", execRenameX, 3)
	RegisterCommand("KEYS", execKeys, 2)
}
