package database

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

// 所有包的init函数都会在main函数之前执行
func init() {
	RegisterCommand("ping", Ping, 1)
}
