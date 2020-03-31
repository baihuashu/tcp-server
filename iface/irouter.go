package iface

//使用了一个模版设计模式 用一个BaseRouter实现这个接口，然后再继承BaseRouter
type IRouter interface {
	PreHandle(request IRequest)
	Handle(request IRequest)
	PostHandle(request IRequest)
}
