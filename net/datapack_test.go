package net

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//测试datapack拆包封包的单元测试
func TestDataPack(t *testing.T){
	/*
	 	模拟的服务器
	 */
	//1 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:",err)
		return
	}
	//创建一个go 承载负责从客户端处理业务
	go func() {
		//2 从客户端读取数据，拆包处理
		for {
			conn,err:= listener.Accept()
			if err != nil {
				fmt.Println("server listen err:",err)
			}

			go func(conn net.Conn) {
				// 拆包的过程
				//定义一个拆包对象
				dp := NewDataPack()
				for{
					//1 第一次从conn读，把包的head读出来
					headData := make([]byte,dp.GetHeadLen())
					io.ReadFull(conn,headData)
					if err!=nil{
						fmt.Println("read head error")
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil{
						fmt.Println("server unpack err",err)
					}
					if msgHead.GetMsgLen()>0 {
						//2 第二次从conn读，根据head中的datalen 再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte,msg.GetMsgLen())
						//根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil{
							fmt.Println("server unpack err",err)
							return
						}
						//完整的一个消息已经读取完毕
						fmt.Println("-->Recv MsgID = ",msg.Id,"datalen = ",msg.DataLen,"data = ",msg.Data)
					}

				}

			}(conn)
		}
	}()

	/*
	模拟客户端
	 */
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err!=nil{
		fmt.Println("client dial err: ",err)
		return
	}
	//创建一个封包对象dp
	dp:=NewDataPack()
	//模拟粘包过程，封装两个msg一同发送
	//封装第一个包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z','i','n','x'},
	}
	sendData1, _ := dp.Pack(msg1)
	//封装第2个包
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'h','e','l','l','o'},
	}
	sendData2, _ := dp.Pack(msg2)
	//将两个包粘到一起
	sendData1 = append(sendData1,sendData2...)
	//一次性发送给服务器
	conn.Write(sendData1)
	//阻塞

	select{}


}