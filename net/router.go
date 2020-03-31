package net

import "github.com/baihuashu/tcp-server/iface"

type BaseRouter struct {
}

/**
 这里BaseRouter的方法都为空，因为有的router不希望有preHandle和PostHandle
 */
func(b *BaseRouter) PreHandle(request iface.IRequest){

}
func(b *BaseRouter) Handle(request iface.IRequest){

}
func(b *BaseRouter) PostHandle(request iface.IRequest){

}
