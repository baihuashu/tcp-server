package main

import (
	"fmt"
	"io"
	"net"
	net2 "tcp-server/net"
	"time"
)

func main() {
	fmt.Println("client start")

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("clinet start err,exit!")
		return
	}

    for {
		//发送封包消息
		dp := net2.NewDataPack()
		bytes, err := dp.Pack(net2.NewMsgPackage(1, []byte("hello word")))
		if err != nil {
			fmt.Println(err)
			return
		}
		if _, err := conn.Write(bytes); err != nil {
			fmt.Println(err)
			return
		}
		//服务器应该给我们回复一个message数据
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, headData); err != nil {
			fmt.Println(err)
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println(err)
		}
		var data []byte
		fmt.Println(msg.GetMsgId(),msg.GetMsgLen())
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println(string(data))
		time.Sleep(time.Second)
	}

}
