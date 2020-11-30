package znet

import (
	"fmt"

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
func newMsgHandle() *MsgHandle {
	MessageHandle = &MsgHandle{
		apis:         make(map[uint32]ziface.IRouter),
		taskQueue:    make([]chan ziface.IRequest, 10),
		workPoolSize: 10,
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
		m.taskQueue[i] = make(chan ziface.IRequest, 50)
		go m.startWork(i)
	}
}

func (m *MsgHandle) startWork(queueID int) {
	fmt.Println("[WORK] work:", queueID, " is starting...")
	queue := m.taskQueue[queueID]

	for {
		select {
		case task, ok := <-queue:
			if ok {
				// fmt.Println("[WORK] work:", queueID, " get task: ", task)
				m.do(task)
			}
		}
	}
}

func init() {
	newMsgHandle()
	MessageHandle.startWorkPool()
}
