package cluster

import (
	"context"
	"errors"
	"go-redis/interface/resp"
	"go-redis/lib/logger"
	"go-redis/lib/utils"
	"go-redis/resp/client"
	"go-redis/resp/reply"
	"strconv"
)

func (c *ClusterDatabase) getPeerClient(peer string) (*client.Client, error) {
	pool, ok := c.nodeConnection[peer]
	if !ok {
		return nil, errors.New("peer not found")
	}
	object, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}
	cli, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("object is not a cli")
	}
	return cli, err
}

func (c *ClusterDatabase) returnPeerClient(peer string, cli *client.Client) error {
	pool, ok := c.nodeConnection[peer]
	if !ok {
		return errors.New("connection not found")
	}
	return pool.ReturnObject(context.Background(), cli)
}

// 转发命令
func (c *ClusterDatabase) relay(peer string, conn resp.Connection, args [][]byte) resp.Reply {
	if peer == c.self {
		return c.db.Exec(conn, args)
	}
	peerClient, err := c.getPeerClient(peer)
	if err != nil {
		return reply.MakeErrReply(err.Error())
	}
	defer func() {
		err := c.returnPeerClient(peer, peerClient)
		if err != nil {
			logger.Error(err)
		}
	}()
	//其他client不知道操作哪个db
	peerClient.Send(utils.ToCmdLine("SELECT", strconv.Itoa(conn.GetDBIndex())))
	return peerClient.Send(args)
}

// 对所有client发送命令 (FLUSHDB)
func (c *ClusterDatabase) broadcast(conn resp.Connection, args [][]byte) map[string]resp.Reply {
	result := make(map[string]resp.Reply)
	for _, peer := range c.nodes {
		result[peer] = c.relay(peer, conn, args)
	}
	return result
}
