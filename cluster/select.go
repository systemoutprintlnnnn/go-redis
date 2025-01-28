package cluster

import "go-redis/interface/resp"

// select指令
func execSelect(c *ClusterDatabase, conn resp.Connection, args [][]byte) resp.Reply {
	return c.db.Exec(conn, args)
}
