package iface

type IServer interface {
	Start()
	Stop()
	Serve()

	AddRouter(msgId uint32,router IRouter)

	GetConnMgr() IConnManager
	//注册OnConnStart 钩子函数的方法
	SetOnConnStart(func(connetction IConnection))
	//注册OnConnStop 钩子函数的方法
	SetOnConnStop(func(connetction IConnection))
	//调用OnConnStart 钩子函数的方法
	CallOnConnStart(connetction IConnection)
	//调用OnConnStop 钩子函数的方法
	CallOnConnStop(connetction IConnection)
}



