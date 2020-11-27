package ziface

import "net"

// IConnection 封装连接
type IConnection interface {
	Start()
	Stop()
	GetConnID() uint32
	GetTCPConnection() *net.TCPConn
	RemoteAddr() net.Addr
	SendMsg(msgID uint32, data []byte) error
	StartRead()
	StartWrite()
}
