package net

import "github.com/baihuashu/tcp-server/iface"

type Request struct {
	conn iface.IConnection

	msg iface.Imessage
}

func (r *Request) GetConnection() iface.IConnection{
	return r.conn
}

func (r *Request) GetData() []byte{
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32{
	return r.msg.GetMsgId()
}