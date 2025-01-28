package cluster

import (
	"go-redis/interface/resp"
)

func newRouter() map[string]CmdFunc {
	router := make(map[string]CmdFunc)
	router["exists"] = defaultFunc // exists key
	router["get"] = defaultFunc    // get key
	router["getset"] = defaultFunc // getset key value
	router["set"] = defaultFunc    // set key value
	router["setnx"] = defaultFunc  // setnx key value
	router["del"] = defaultFunc    // del key
	router["ping"] = ping
	router["rename"] = rename
	router["renamenx"] = rename
	router["flushdb"] = flushdb
	router["del"] = del
	router["select"] = execSelect

	return router
}

// 单纯转发 GET SET DEL
func defaultFunc(c *ClusterDatabase, conn resp.Connection, args [][]byte) resp.Reply {
	key := string(args[1])
	peer := c.nodeSelector.GetNode(key)
	return c.relay(peer, conn, args)
}
