package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// rename k1 k2
func rename(c *ClusterDatabase, conn resp.Connection, args [][]byte) resp.Reply {
	if len(args) != 3 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'rename' command")
	}
	//判断k1是否存在
	k1 := string(args[1])
	k2 := string(args[2])
	n1 := c.nodeSelector.GetNode(k1)
	n2 := c.nodeSelector.GetNode(k2)
	if n1 != n2 {
		return reply.MakeErrReply("ERR keys in different slots")
	}

	return c.relay(n1, conn, args)
}
