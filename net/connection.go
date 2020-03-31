package net

import (
	"errors"
	"fmt"
	"io"
	"net"
	"github.com/baihuashu/tcp-server/iface"
	"github.com/baihuashu/tcp-server/utils"
)

type Connection struct {
	//为了在Connection方法中 添加conn到connMgr
	TcpServer iface.IServer

	//当前链接
	Conn *net.TCPConn

	ConnID uint32

	isClosed bool
	//无缓冲的管道，用于读,写groutine之间的消息通信
	msgChan chan []byte
	//告知当前链接退出的channel(由Reader告知Writer退出)
	ExitChan chan bool
	//消息的管理MsgID 和对应的处理业务API关系 （从server那继承过来）
	MsgHandler iface.IMsgHandle
}

func NewConnection(server iface.IServer, conn *net.TCPConn, connID uint32, msgHandler iface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool, 1),
	}
	//将conn 加入到ConnMgr中
	//todo 但是Connection中并没有ConnMgr属性 要怎么办呢？[在Connection中添加server的字段 指回去]
	c.TcpServer.GetConnMgr().Add(c)
	return c
}
func (c *Connection) StartReader() {
	fmt.Println("Reader goroutine is running ...")
	defer fmt.Println("connId = ", c.ConnID, "[Reader is exit] ,remote addr is ", c.RemoteAddr())
	defer c.Stop()
	for {
		//buf:= make([]byte,utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err !=nil{
		//	fmt.Println("rev buf err",err)
		//	continue
		//}
		dp := NewDataPack()
		//读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error", err)
			break
		}
		//拆包，得到MsgId和msgDatalen 放在msg中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		//根据datalen，再次读取data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg head error", err)
				break
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息发送给Worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//这里 有多少个连接 就有多少个协程 进行升级，交给工作池
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

/*
	写消息Groutine,专门发送给客户端消息的模块
*/
func (c *Connection) StartWritr() {
	fmt.Println("[Writer Groutine is running....]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")
	//不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error", err)
				return
			}
		//Reader出意外，就终止Writer
		//这里有个技巧是： 看StartReader中如果出现异常就会break 最终调用stop方法，在stop方法中关闭
		case <-c.ExitChan:
			return
		}
	}
}

//提供一个SendMsg方法 将我们要发送给客户端的数据，先进性封包，再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send msg")
	}
	//将data进行封包 MsgDataLen/MsgId/Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("xixi", string(binaryMsg))

	//通过channel让写业务执行  将数据发送给客户端
	c.msgChan <- binaryMsg
	return nil
}

//让当前链接开始工作
func (c *Connection) Start() {
	fmt.Println("conn start()... ConnID = ", c.ConnID)
	go c.StartReader()
	go c.StartWritr()

	//创建conn后 调用对应Hookef
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("conn stop ConnID=", c.ConnID)
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//调用 销毁链接之前的hook
	c.TcpServer.CallOnConnStart(c)

	//关闭socket链接
	c.Conn.Close()
	//告知Writer关闭
	c.ExitChan <- true

	//将当前连接从ConnMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源？
	close(c.ExitChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的ip，port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
