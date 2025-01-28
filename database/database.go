package database

import (
	"go-redis/aof"
	"go-redis/config"
	"go-redis/interface/resp"
	"go-redis/lib/logger"
	"go-redis/resp/reply"
	"strconv"
	"strings"
)

type StandaloneDatabase struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
}

func NewStandaloneDatabase() *StandaloneDatabase {
	database := &StandaloneDatabase{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}

	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := 0; i < config.Properties.Databases; i++ {
		db := NewDB()
		db.index = i
		database.dbSet[i] = db
	}
	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAOFHandler(database)
		if err != nil {
			logger.Error("failed to create aof handler: %v", err)
			panic(err)
		}
		database.aofHandler = aofHandler
		for _, db := range database.dbSet {
			db.addAof = func(line CmdLine) {
				database.aofHandler.AddAof(db.index, line)
			}
		}

	}
	return database
}

// SET k v
// GET k
// SELECT 2
func (d *StandaloneDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply("ERR wrong number of arguments for 'select' command")
		}
		return execSelect(client, d, args[1:])
	}
	dbIndex := client.GetDBIndex()
	db := d.dbSet[dbIndex]
	return db.Exec(client, args)

}

func (d StandaloneDatabase) Close() {
	//TODO implement me
	panic("implement me")
}

func (d StandaloneDatabase) AfterClientClose(client resp.Connection) {
	//TODO implement me
	panic("implement me")
}

// SELECT 2
func execSelect(c resp.Connection, database *StandaloneDatabase, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}
	if dbIndex < 0 || dbIndex >= len(database.dbSet) {
		return reply.MakeErrReply("ERR DB index is out of range")
	}

	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}
