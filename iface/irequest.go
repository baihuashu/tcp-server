package iface

//需要把客户端请求的 连接信息和请求数据 ，封装到Request。
type IRequest interface {
	GetConnection() IConnection

	GetData() []byte

	GetMsgID() uint32
}
