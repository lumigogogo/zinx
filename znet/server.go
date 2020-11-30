package znet

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lumigogogo/zinx/ziface"
)

// Server 监听连接
type Server struct {
	Name       string         // server name
	IPVersion  string         // ipv4
	IP         string         // listen ip
	Port       int            // listen port
	MsgHandler *MsgHandle     // control router
	ExitChan   chan os.Signal // exit channel
}

var connID uint32

// Start 启动服务
func (s *Server) Start() {
	fmt.Println("[START] server name: ", s.Name, " listen ip: ", s.IP, " port: ", s.Port, " is start!")

	go func() {
		listenAddr := fmt.Sprintf("%s:%d", s.IP, s.Port)
		addr, err := net.ResolveTCPAddr(s.IPVersion, listenAddr)
		if err != nil {
			fmt.Println("[START] create addr: ", listenAddr, " error!")
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("[START] listen addr: ", listenAddr, " error! ", err)
			return
		}

		fmt.Println("[START] server is listening!")

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("[START] listener accept error! err: ", err)
				continue
			}
			// fmt.Println("[START] get conn remote: ", conn.RemoteAddr())

			atomic.AddUint32(&connID, 1)
			dealConn := NewConnection(connID, conn, s)
			fmt.Println("[START] get conn id, ", connID)

			go dealConn.Start()
		}
	}()
}

// Stop 停止服务
func (s *Server) Stop() {

}

// Serve 监听函数
func (s *Server) Serve() {
	s.Start()

	signal.Notify(s.ExitChan, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	<-s.ExitChan
	fmt.Println("[Serve] server will be stop after 3 second!")
	time.Sleep(3 * time.Second)
}

// AddRouter add router to Server.apis
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.apis[msgID] = router
}

// NewServer create new server
func NewServer(name, ip string, port int) ziface.IServer {
	return &Server{
		Name:       name,
		IPVersion:  "tcp4",
		IP:         ip,
		Port:       port,
		MsgHandler: MessageHandle,
		ExitChan:   make(chan os.Signal),
	}
}
