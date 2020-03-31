package main

import (
	"fmt"
	"tcp-server/iface"
	"tcp-server/net"
)

type PingRouter struct {
	net.BaseRouter
}


func(this *PingRouter) Handle(request iface.IRequest){
	fmt.Println("Call Router Handle")
	//先读数据 再回写ping ping ping
	fmt.Println("recv from client : msgId= ",request.GetMsgID(),"data =",string(request.GetData()))
	if err := request.GetConnection().SendMsg(1, []byte("ping..ping..ping"));err != nil{
		fmt.Println("Handle error:", err)
	}
}

type HelloRouter struct {
	net.BaseRouter
}


func(this *HelloRouter) Handle(request iface.IRequest){
	fmt.Println("Call Router Handle")
	//先读数据 再回写ping ping ping
	fmt.Println("recv from client : msgId= ",request.GetMsgID(),"data =",string(request.GetData()))
	if err := request.GetConnection().SendMsg(1, []byte("hello..hello..hello"));err != nil{
		fmt.Println("Handle error:", err)
	}
}

//conn创建 的钩子函数
func DoConnectionBegin(conn iface.IConnection){
	fmt.Println("===> DoConnectionLost is Called")
	if err := conn.SendMsg(202,[]byte("DoConnection begin")); err !=nil {
		fmt.Println(err)
	}
}
//conn断开前的hook函数
func DoConnectionLost(conn iface.IConnection){
	fmt.Println("===> DoConnectionBegin is Called")
	fmt.Println("conn ID = ",conn.GetConnID(),"is Lost")
}

func main()  {
	server := net.NewServer("wzl-server")
	server.AddRouter(0,&PingRouter{})
	server.AddRouter(1,&HelloRouter{})
	server.SetOnConnStart(DoConnectionBegin)
	server.SetOnConnStop(DoConnectionLost)
	server.Serve()
}