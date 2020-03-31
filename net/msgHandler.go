package net

import (
	"fmt"
	"strconv"
	"github.com/baihuashu/tcp-server/iface"
	"github.com/baihuashu/tcp-server/utils"
)

type MsgHandle struct {
	//存放每个MsgID所对应的处理方法
	Apis map[uint32] iface.IRouter

	//负责Worker取任务的消息队列
	TaskQueue []chan iface.IRequest
	//业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

func NewMsgHandle() *MsgHandle{
	return &MsgHandle{
		Apis:make(map[uint32] iface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue: make([]chan iface.IRequest,utils.GlobalObject.WorkerPoolSize),
	}
}


func (mh *MsgHandle) DoMsgHandler(request iface.IRequest){
	//1 从request中找到msgId
	handler,ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID =",request.GetMsgID(),"is not found! Need Register!")
	}
	//2 根据MsgID调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}
func (mh *MsgHandle) AddRouter(msgID uint32,router iface.IRouter){
	if _,ok  := mh.Apis[msgID]; ok {
		//id 已经注册
		panic("repeat api, msgID = "+strconv.Itoa(int(msgID)))
	}
	mh.Apis[msgID]=router
	fmt.Println("Add api MsgId",msgID,"succ!")
}

//启动一个Worker池 (开启工作吃的动作只能发生一次，一个框架只能有一个worker工作池)
func (mh *MsgHandle) StartWorkerPool(){
	//根据workerPoolSize 分别开启worker，每个worker用一个go来承载
	for i:=0;i<int(mh.WorkerPoolSize);i++{
		//一个worker被启动
		//1 当前的worker对应的channel消息队列 开辟空间 第0个worker 就用第0个channel ...
		mh.TaskQueue[i] = make(chan iface.IRequest,utils.GlobalObject.MaxWorkerTaskLen)
		//2 启动当前的worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i,mh.TaskQueue[i])
	}
}

//启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int,taskQueue chan iface.IRequest){
	fmt.Println("Worker ID = ",workerID,"is started ..")
	//不断的阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的Request,执行当前Request所绑定的业务
		case request:= <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}
//将消息交给TaskQueue，由Worker来处理
func (mh *MsgHandle)SendMsgToTaskQueue(request iface.IRequest){
	//1 将消息平均分配给不通过的worker
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	//2 将消息发送 给对应的worker的TaskQueue即可
	fmt.Printf("Add ConnId = %d ,request MsgID = %d, to WorkerID = %d",request.GetConnection().GetConnID(),request.GetMsgID(),workerID)
	mh.TaskQueue[workerID] <- request
}
