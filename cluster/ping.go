package cluster

import "go-redis/interface/resp"

// ping指令
func ping(c *ClusterDatabase, conn resp.Connection, args [][]byte) resp.Reply {
	return c.db.Exec(conn, args)
}
