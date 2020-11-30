package znet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/lumigogogo/zinx/ziface"
)

// Connection ..
type Connection struct {
	connID   uint32
	conn     *net.TCPConn
	s        ziface.IServer
	msgChan  chan []byte
	isClosed bool
	// sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// Start do ...
func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.StartRead()
	go c.StartWrite()
}

// Stop close
func (c *Connection) Stop() {
	fmt.Println("[Stop] conn: ", c.connID, " is stoping...")
	// c.Lock()
	// defer c.Unlock()

	if c.isClosed {
		return
	}
	c.isClosed = true

	c.conn.Close()
	c.cancel()
	close(c.msgChan)
}

// GetConnID ..
func (c *Connection) GetConnID() uint32 {
	return c.connID
}

// GetTCPConnection ..
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.conn
}

// RemoteAddr ..
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// SendMsg ..
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	msg := NewMessage(msgID, uint32(len(data)), data)

	b, err := Pack(msg)
	if err != nil {
		return errors.New("[SendMsg] pack msg error")
	}
	c.msgChan <- b
	return nil
}

// StartRead ..
func (c *Connection) StartRead() {
	defer c.Stop()

	select {
	case <-c.ctx.Done():
		return
	default:
		for {
			headData := make([]byte, DataHeadLen)
			_, err := io.ReadFull(c.conn, headData) //ReadFull 会把msg填充满为止
			if err != nil {
				fmt.Println("read head error: ", err)
				return
			}
			msgHead, err := Unpack(headData)
			if err != nil {
				fmt.Println("server unpack err:", err)
				return
			}

			if msgHead.GetDataLen() > 0 {
				msg := msgHead.(*Message)
				msg.SetData(make([]byte, msg.GetDataLen()))

				//根据dataLen从io中读取字节流
				_, err := io.ReadFull(c.conn, msg.GetData())
				if err != nil {
					fmt.Println("server unpack data err:", err)
					return
				}

				request := &Request{
					conn: c,
					msg:  msg,
				}

				c.s.(*Server).MsgHandler.SendTaskToTaskQueue(request)
			}
		}
	}
	fmt.Println("[StartRead] wait done!")
}

// StartWrite ..
func (c *Connection) StartWrite() {
	defer fmt.Println("[StartWrite] conn id: ", c.connID, " write is done!")
	for {
		select {
		case <-c.ctx.Done():
			fmt.Println("[StartWrite] ctx done")
			return

		case data, ok := <-c.msgChan:
			if ok {
				_, err := c.conn.Write(data)
				if _, ok := err.(*net.OpError); ok {
					return
				}
				if err != nil {
					fmt.Println("[StartWrite] write to client error, ", err)
				}
			}
		}
	}
}

// NewConnection return new conn
func NewConnection(connID uint32, conn *net.TCPConn, server ziface.IServer) ziface.IConnection {
	return &Connection{
		connID:   connID,
		conn:     conn,
		s:        server,
		isClosed: false,
		msgChan:  make(chan []byte, 512),
	}
}
