package database

import "go-redis/interface/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	Close()
	AfterClientClose(client resp.Connection)
}

//	type Database interface {
//		Exec(client resp.Connection, args [][]byte) resp.Reply
//		AfterClientClose(c resp.Connection)
//		Close()
//	}
type DataEntity struct {
	Data interface{}
}

func NewDataEntity(data interface{}) *DataEntity {
	return &DataEntity{
		Data: data,
	}
}
