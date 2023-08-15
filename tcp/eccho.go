package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/sunflower10086/go-redis/lib/logger"
	"github.com/sunflower10086/go-redis/lib/sync/atomic"
	"github.com/sunflower10086/go-redis/lib/sync/wait"
)

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	e.Conn.Close()
	return nil
}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (e *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if e.closing.Get() {
		conn.Close()
	}

	client := EchoClient{
		Conn: conn,
	}

	e.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				e.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}

			return
		}

		client.Waiting.Add(1)
		conn.Write([]byte(msg))
		client.Waiting.Done()
	}
}

func (e *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	e.closing.Set(true)

	e.activeConn.Range(func(key, value any) bool {
		key.(EchoClient).Conn.Close()
		return true
	})

	return nil
}
