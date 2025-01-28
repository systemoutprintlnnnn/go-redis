package tcp

import (
	"bufio"
	"context"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

type EchoHandler struct {
	//有效连接数
	activeConn sync.Map
	//是否正在关闭
	closing atomic.Boolean
}

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

func (e EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

func (handler *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	handler.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		//The line reads a string from the reader until a newline character ('\n') is encountered.
		//It assigns the read string to msg and any error encountered during the read operation to err.
		msg, err := reader.ReadString('\n')
		if err != nil {
			//操作系统中数据结束符
			if err == io.EOF {
				logger.Info("Connection close")
				handler.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		logger.Info("msg:", msg)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}

}

func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	//状态改为正在关闭
	handler.closing.Set(true)
	//客户端全部关掉
	//遍历sync.map的方法
	handler.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Conn.Close()
		return true
	})
	return nil
}
