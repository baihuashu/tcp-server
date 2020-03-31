package iface

import "net"

//将每个conn封装。
//疑问：为啥*net.TCPConn需要指针但是net.Addr不需要
type IConnection interface {
	Start()

	Stop()

	GetTCPConnection() *net.TCPConn

	GetConnID() uint32
	//获取远程客户端的ip，port
	RemoteAddr() net.Addr
	//发送数据给远程客户端
	SendMsg(msgId uint32,data []byte) error
}



