package net

import (
	"fmt"
	"net"
	"github.com/baihuashu/tcp-server/iface"
	"github.com/baihuashu/tcp-server/utils"
)

type Server struct {
	Name string
	//服务器绑定的ip版本
	IPVersion string
	IP        string
	Port      int
	//当前server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler iface.IMsgHandle
	//该server的连接管理器
	ConnMgr iface.IConnManager
	//该Server创建链接之后自动调用Hook函数 --OnConnStart
	OnConnStart func(conn iface.IConnection)
	//该Server销毁链接之前自动调用Hook函数 --OnConnStop
	OnConnStop func(conn iface.IConnection)
}

func (s *Server) Start() {
	go func() {
		//0 开启工作池
		s.MsgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		fmt.Printf("[Start] Server Listenner at IP: %s ,Port %d ,is starting\n", s.IP, s.Port)
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println(err)
			return
		}
		//2。监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[Success] server: %s is listening... \n", s.Name)
		//3。阻塞的等待客户端链接，处理客户端链接业务（读写）
		var cid uint32
		cid = 0
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println(err)
				continue
			}
			//设置最大连接个数的判断，如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				//todo 给客户端相应一个超出最大连接的错误包
				fmt.Println("too many conn !")
				conn.Close()
			}

			//与客户端建立链接，做一些业务。
			dealConn := NewConnection(s,conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	//todo 将一些服务器的资源，状态或者一些开辟的链接信息 进行停止或者回收
	fmt.Println("[STOP] server is stop")
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	select {}
}
func (s *Server) AddRouter(msgId uint32, router iface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router Success")
}
//看 connection中的new方法
func (s *Server) GetConnMgr() iface.IConnManager {
	return s.ConnMgr
}

func NewServer(name string) iface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnMgr(),
	}
	return s
}

//注册OnConnStart 钩子函数的方法
func(s *Server) SetOnConnStart(hookFunc func(connetction iface.IConnection)){
	s.OnConnStart = hookFunc
}
//注册OnConnStop 钩子函数的方法
func(s *Server) SetOnConnStop(hookFunc func(connetction iface.IConnection)){
	s.OnConnStop = hookFunc
}
//调用OnConnStart 钩子函数的方法
func(s *Server) CallOnConnStart(conn iface.IConnection){
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStart()")
		s.OnConnStart(conn)
	}
}
//调用OnConnStop 钩子函数的方法
func(s *Server) CallOnConnStop(conn iface.IConnection){
	if s.OnConnStop != nil {
		fmt.Println("---> Call OnConnStop()")
		s.OnConnStop(conn)
	}
}
