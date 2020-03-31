package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/baihuashu/tcp-server/iface"
	"github.com/baihuashu/tcp-server/utils"
)

//封包，拆包的具体模块
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	//DataLen unit32(4字节)+ID unit32(4字节)
	return 8
}
func (dp *DataPack) Pack(msg iface.Imessage) ([]byte, error) {

	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//将dataLen写进buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将MsgId写进buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将data写进buff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//拆包方法：（将包的Head信息读出来）   后续再根据head信息里的data长度，再进行一次读（此方法不含有这个过程）
func (dp *DataPack) UnPack(binaryData []byte) (iface.Imessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)
	//只解压head信息，得到datalen和MsgID
	msg := &Message{}
	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//先判断 datalen是否超出我们允许的最大长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	//if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen);err !=nil{
	//	return nil ,err
	//}
	return msg, nil

}
