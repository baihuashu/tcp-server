package iface


/*
	消息管理抽象层
	为多路由服务
 */

type IMsgHandle interface {
	//调度/执行响应的Router消息处理方法
	DoMsgHandler(request IRequest)
	//为消息添加具体的处理逻辑
	AddRouter(msgID uint32,router IRouter)

	StartWorkerPool()

	SendMsgToTaskQueue(request IRequest)
}