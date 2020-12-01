package znet

import "github.com/lumigogogo/zinx/ziface"

// Request ..
type Request struct {
	conn ziface.IConnection
	msg  ziface.IMessage
}

// GetConnection ..
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// GetData ..
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID ..
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
