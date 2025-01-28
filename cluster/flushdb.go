package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

func flushdb(c *ClusterDatabase, conn resp.Connection, args [][]byte) resp.Reply {
	replies := c.broadcast(conn, args)
	for _, r := range replies {
		if reply.IsErrorReply(r) {
			return r
		}
	}
	return reply.MakeOkReply()
}
