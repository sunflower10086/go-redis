package handler

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"

	"github.com/sunflower10086/go-redis/database"
	databasefase "github.com/sunflower10086/go-redis/interface/database"
	"github.com/sunflower10086/go-redis/lib/logger"
	"github.com/sunflower10086/go-redis/lib/sync/atomic"
	"github.com/sunflower10086/go-redis/resp/connection"
	"github.com/sunflower10086/go-redis/resp/parser"
	"github.com/sunflower10086/go-redis/resp/reply"
)

type RespHandler struct {
	activeConn sync.Map
	db         databasefase.Database
	closing    atomic.Boolean
}

func NewRespHandler() *RespHandler {
	var db databasefase.Database
	db = database.NewEchoDatabase()
	//TODO: 实现database
	return &RespHandler{
		db: db,
	}
}

// CloseClient 关闭其中一个连接
func (r *RespHandler) CloseClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}
	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})

	ch := parser.ParseStream(conn)
	for payload := range ch {
		// 有错误
		if payload.Error != nil {
			err := payload.Error
			if err == io.EOF ||
				err == io.ErrUnexpectedEOF ||
				errors.Is(err, errors.New("use of close network connection")) {
				r.CloseClient(client)
				logger.Info("connection close: " + client.RemoteAddr().String())
				return
			}
			// protocol error
			errReply := reply.NewErrReply(err.Error())
			if err := client.Write(errReply.ToBytes()); err != nil {
				r.CloseClient(client)
				logger.Info("connection close: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		// Exec
		if payload.Data == nil {
			continue
		}

		data, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			continue
		}
		result := r.db.Exec(client, data.Args)
		if result != nil {
			client.Write(result.ToBytes())
		} else {
			client.Write(reply.UnknownErrBytes)
		}
	}
}

// Close 关闭所有连接
func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(func(key, value any) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})

	r.db.Close()
	return nil
}
