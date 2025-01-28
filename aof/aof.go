package aof

import (
	"go-redis/config"
	databaseface "go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/utils"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"io"
	"os"
	"strconv"
)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

const (
	aofQueueSize = 1 << 16
)

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

// AofHandler receive msgs from channel and write to AOF file
type AofHandler struct {
	db          databaseface.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	//记录上一条命令的数据库索引，用于判断是否需要切换数据库
	currentDB int
}

// NewAOFHandler creates a new aof.AofHandler
func NewAOFHandler(db databaseface.Database) (*AofHandler, error) {
	handler := &AofHandler{}
	handler.aofFilename = config.Properties.AppendFilename
	if config.Properties.AppendFilename == "" && config.Properties.AppendOnly {
		panic("appendfilename is empty")
	}
	handler.db = db
	//加载AOF文件
	handler.LoadAof()
	aofFile, err := os.OpenFile(handler.aofFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofFile
	handler.aofChan = make(chan *payload, aofQueueSize)
	go func() {
		handler.handleAof()
	}()
	return handler, nil
}

func (handler *AofHandler) AddAof(dbIndex int, cmd CmdLine) {
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmd,
			dbIndex: dbIndex,
		}

	}
}

// 存储到AOF文件
func (handler *AofHandler) handleAof() {
	handler.currentDB = 0
	for p := range handler.aofChan {
		if p.dbIndex != handler.currentDB {
			data := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				continue
			}
			handler.currentDB = p.dbIndex
		}

		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Error(err)
			continue
		}
	}

}

// 读取AOF文件
func (handler *AofHandler) LoadAof() {
	file, err := os.Open(handler.aofFilename)
	if err != nil {
		logger.Error(err)
		return
	}
	defer file.Close()

	ch := parser.ParseStream(file)
	fakeConn := connection.NewConnection(nil)
	for payload := range ch {
		//ch得到错误
		if payload.Err != nil {
			//文件读完了
			if payload.Err == io.EOF {
				break
			}
			logger.Error(payload.Err)
			continue
		}
		//ch得到正常数据
		if payload.Data == nil {
			logger.Error("nil payload data")
			continue
		}

		r, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("invalid reply type")
			continue
		}

		rep := handler.db.Exec(fakeConn, r.Args)
		if reply.IsErrorReply(rep) {
			logger.Error("AOF exec error: %s", rep.ToBytes())
		}

	}
}

//// LoadAof read aof file
//func (handler *AofHandler) LoadAof() {
//
//	file, err := os.Open(handler.aofFilename)
//	if err != nil {
//		logger.Warn(err)
//		return
//	}
//	defer file.Close()
//	ch := parser.ParseStream(file)
//	fakeConn := &connection.Connection{} // only used for save dbIndex
//	for p := range ch {
//		if p.Err != nil {
//			if p.Err == io.EOF {
//				break
//			}
//			logger.Error("parse error: " + p.Err.Error())
//			continue
//		}
//		if p.Data == nil {
//			logger.Error("empty payload")
//			continue
//		}
//		r, ok := p.Data.(*reply.MultiBulkReply)
//		if !ok {
//			logger.Error("require multi bulk reply")
//			continue
//		}
//		ret := handler.db.Exec(fakeConn, r.Args)
//		if reply.IsErrorReply(ret) {
//			logger.Error("exec err", err)
//		}
//	}
//}
