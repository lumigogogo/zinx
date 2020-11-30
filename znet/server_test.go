package znet

import (
	"testing"

	"github.com/lumigogogo/zinx/ziface"
)

type PingRouter struct {
	Router
}

func (p *PingRouter) PreHandle(request ziface.IRequest) {

}

func (p *PingRouter) Handle(request ziface.IRequest) {
	request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
}

func (p *PingRouter) AfterHandle(request ziface.IRequest) {

}

type PangRouter struct {
	Router
}

func (pr *PangRouter) PreHandle(request ziface.IRequest) {

}

func (pr *PangRouter) Handle(request ziface.IRequest) {
	request.GetConnection().SendMsg(1, []byte("pang...pang...pang"))
}

func (pr *PangRouter) AfterHandle(request ziface.IRequest) {

}

func TestRouter(t *testing.T) {
	go func() {
		s := NewServer("test-router", "127.0.0.1", 8888)
		s.AddRouter(1, &PingRouter{})
		s.AddRouter(2, &PangRouter{})
		s.Serve()
	}()
}
