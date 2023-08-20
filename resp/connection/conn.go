package connection

import (
	"net"
	"sync"
	"time"

	"github.com/sunflower10086/go-redis/lib/sync/wait"
)

// Connection 描述客户的结构体
type Connection struct {
	Conn     net.Conn
	Waiting  wait.Wait
	mu       sync.Mutex
	selectDB int
}

func NewConn(conn net.Conn) *Connection {
	return &Connection{
		Conn: conn,
	}
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	_ = c.Conn.Close()
	return nil
}

func (c *Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	// 同时只能有一个人去写客户端
	c.mu.Lock()
	c.Waiting.Add(1)
	defer func() {
		c.mu.Unlock()
		c.Waiting.Done()
	}()
	_, err := c.Conn.Write(bytes)
	return err
}

func (c *Connection) GetDBIndex() int {
	return c.selectDB
}

func (c *Connection) SelectDB(i int) {
	//TODO implement me
	panic("implement me")
}
