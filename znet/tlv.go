package znet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/lumigogogo/zinx/utils"
	"github.com/lumigogogo/zinx/ziface"
)

// DataHeadLen ..
const DataHeadLen = 8

// Pack 按照TLV格式封装数据, TLV: tag|type|data_len|data
func Pack(msg ziface.IMessage) ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})

	if err := binary.Write(buff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	if err := binary.Write(buff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(buff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// Unpack 解码TLV并封装成IMessage
func Unpack(data []byte) (ziface.IMessage, error) {
	reader := bytes.NewReader(data)

	msg := &Message{}

	if err := binary.Read(reader, binary.LittleEndian, &msg.MsgID); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if utils.GlobalConf.MaxPacketSize > 0 && msg.DataLen > utils.GlobalConf.MaxPacketSize {
		return nil, errors.New("too large msg data received")
	}

	return msg, nil
}
