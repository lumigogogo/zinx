package znet

import "github.com/lumigogogo/zinx/ziface"

// Message ..
type Message struct {
	MsgID   uint32
	DataLen uint32
	data    []byte
}

// NewMessage ..
func NewMessage(msgID, dataLen uint32, data []byte) ziface.IMessage {
	return &Message{
		MsgID:   msgID,
		DataLen: dataLen,
		data:    data,
	}
}

// GetDataLen ..
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// GetMsgID ..
func (m *Message) GetMsgID() uint32 {
	return m.MsgID
}

// GetData ..
func (m *Message) GetData() []byte {
	return m.data
}

// SetDataLen ..
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

// SetMsgID ..
func (m *Message) SetMsgID(id uint32) {
	m.MsgID = id
}

// SetData ..
func (m *Message) SetData(data []byte) {
	m.data = data
}
