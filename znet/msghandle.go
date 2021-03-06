package znet

import (
	"fmt"

	"github.com/lumigogogo/zinx/utils"
	"github.com/lumigogogo/zinx/ziface"
)

// MsgHandle 根据业务ID执行具体的router
type MsgHandle struct {
	apis         map[uint32]ziface.IRouter
	workPoolSize int
	taskQueue    []chan ziface.IRequest
}

// MessageHandle ..
var MessageHandle *MsgHandle

// NewMsgHandle create new msg handle..
func NewMsgHandle() *MsgHandle {
	MessageHandle = &MsgHandle{
		apis:         make(map[uint32]ziface.IRouter),
		taskQueue:    make([]chan ziface.IRequest, utils.GlobalConf.WorkPoolSize),
		workPoolSize: utils.GlobalConf.WorkPoolSize,
	}
	return MessageHandle
}

// SendTaskToTaskQueue ..
func (m *MsgHandle) SendTaskToTaskQueue(request ziface.IRequest) {
	queueID := (int)(request.GetConnection().GetConnID()) % m.workPoolSize
	m.taskQueue[queueID] <- request
}

// Do ..
func (m *MsgHandle) do(request ziface.IRequest) {
	router := m.apis[request.GetMsgID()]

	router.AfterHandle(request)
	router.Handle(request)
	router.AfterHandle(request)
}

func (m *MsgHandle) startWorkPool() {
	for i := 0; i < m.workPoolSize; i++ {
		m.taskQueue[i] = make(chan ziface.IRequest, utils.GlobalConf.MaxWorkTaskNum)
		go m.startWork(i)
	}
	fmt.Println("[WORK] work num:", m.workPoolSize, " is starting")
}

func (m *MsgHandle) startWork(queueID int) {
	queue := m.taskQueue[queueID]

	for {
		select {
		case task, ok := <-queue:
			if ok {
				m.do(task)
			}
		}
	}
}
