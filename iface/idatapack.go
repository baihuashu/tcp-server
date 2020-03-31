package iface

/**
	封包，拆包
	直接面向TCP连接中的数据流，用于处理TCP粘包问题
 */
type IdataPack interface {

	GetHeadLen() uint32
	//封包方法
	Pack(msg Imessage) ([]byte,error)
	//拆包方法
	UnPack([]byte)(Imessage,error)



}
