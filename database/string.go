package database

import (
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/lib/utils"
	"go-redis/resp/reply"
)

// GET k1
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	//判断entity的类型
	if t, ok := entity.Data.([]byte); ok {
		return reply.MakeBulkReply(t)
	} else {
		return reply.MakeErrReply("value is not a string")
	}

}

//SET k v

func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	db.PutEntity(key, database.NewDataEntity(value))
	db.addAof(utils.ToCmdLine2("SET", args...))
	return reply.MakeOkReply()
}

// SETNX k1 v1
// SET if Not exists
func execSetnx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	db.addAof(utils.ToCmdLine2("SETNX", args...))
	return reply.MakeIntReply(int64(db.PutIfAbsent(key, database.NewDataEntity(value))))
}

// GETSET k v
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity, ok := db.GetEntity(key)
	db.PutEntity(key, database.NewDataEntity(value))
	db.addAof(utils.ToCmdLine2("GETSET", args...))
	if !ok {
		return reply.MakeNullBulkReply()
	}
	//判断entity的类型
	if t, ok := entity.Data.([]byte); ok {
		return reply.MakeBulkReply(t)
	} else {
		return reply.MakeErrReply("value is not a string")
	}
}

// STRLEN k
func execStrlen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, ok := db.GetEntity(key)
	if !ok {
		return reply.MakeNullBulkReply()
	}
	//判断entity的类型
	if t, ok := entity.Data.([]byte); ok {
		return reply.MakeIntReply(int64(len(t)))
	} else {
		return reply.MakeErrReply("value is not a string")
	}
}

func init() {
	RegisterCommand("get", execGet, 2)
	RegisterCommand("set", execSet, 3)
	RegisterCommand("setnx", execSetnx, 3)
	RegisterCommand("getset", execGetSet, 3)
	RegisterCommand("strlen", execStrlen, 2)
}
