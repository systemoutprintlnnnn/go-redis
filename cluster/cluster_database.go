package cluster

import (
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"go-redis/config"
	database2 "go-redis/database"
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/lib/consistenthash"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
	"strings"
)

type ClusterDatabase struct {
	self string
	// 保存所有节点的地址
	nodes []string
	//节点选择器
	nodeSelector *consistenthash.NodeMap
	//保存所有节点的连接池
	nodeConnection map[string]*pool.ObjectPool
	//保存所有standalone数据库
	db database.Database
}

var router = newRouter()

func NewClusterDatabase() *ClusterDatabase {
	c := &ClusterDatabase{
		self:           config.Properties.Self,
		db:             database2.NewStandaloneDatabase(),
		nodeSelector:   consistenthash.NewNodeMap(nil),
		nodeConnection: make(map[string]*pool.ObjectPool),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, c.self)
	c.nodeSelector.AddNode(nodes...)
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		c.nodeConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{node: peer})
	}
	c.nodes = nodes
	return c
}

type CmdFunc func(c *ClusterDatabase, conn resp.Connection, args [][]byte) resp.Reply

func (c *ClusterDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			res := reply.MakeUnknownErrReply()
			logger.Error(err)
			logger.Error(res)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return reply.MakeErrReply("unknown command '" + cmdName + "'")
	}
	return cmdFunc(c, client, args)
}

func (c *ClusterDatabase) Close() {
	c.Close()
}

func (c *ClusterDatabase) AfterClientClose(client resp.Connection) {
	c.AfterClientClose(client)
}
