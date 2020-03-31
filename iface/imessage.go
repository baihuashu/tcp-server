package iface

type Imessage interface {
	//消息id
	GetMsgId() uint32
	//消息长度
	GetMsgLen() uint32
	//消息内容
	GetData() []byte

	SetMsgId(uint32)
	SetMsgLen(uint32)
	SetData([]byte)
}
