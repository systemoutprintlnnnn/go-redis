package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// del k1 k2 k3...
func del(c *ClusterDatabase, conn resp.Connection, args [][]byte) resp.Reply {
	replies := c.broadcast(conn, args)
	var deleted int64 = 0
	for _, r := range replies {
		if reply.IsErrorReply(r) {
			return r
		}
		intReply, ok := r.(*reply.IntReply)
		if !ok {
			return reply.MakeErrReply("ERR wrong type")
		}
		deleted += intReply.Code
	}

	return reply.MakeIntReply(deleted)
}
