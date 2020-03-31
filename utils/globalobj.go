package utils

import (
	"encoding/json"
	"io/ioutil"
	"github.com/baihuashu/tcp-server/iface"
)

/**
 	储存全局参数，供其他模块使用
 */
type GlobalObj struct {
	TcpServer iface.IServer //当前全局Server对象
	Host string             //当前服务器监听ip
	TcpPort int 			//当前服务器监听端口号
	Name string 			//当前服务器名称

	Version string		    //当前版本号
	MaxConn int 			//当前服务器允许的最大连接数
	MaxPackageSize uint32   //当前框架数据包的最大值

	WorkerPoolSize uint32 	//当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32	//框架允许用户最多开辟多少个Worker（限定条件）
}
/**
	定义一个全局对外的GlobalObj
 */
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload(){
	bytes, err := ioutil.ReadFile("conf/server.json")
	if err !=nil{
		//没必要执行下去了
		panic(err)
	}
	//将json文件数据解析到struct中
	err = json.Unmarshal(bytes, &GlobalObject)
	if err != nil{
		panic(err)
	}
}
/**
	初始化方法。
 */
func init(){
	//如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		TcpServer:      nil,
		Host:           "0.0.0.0",
		TcpPort:        8999,
		Name:           "wzlServerApp",
		Version:        "",
		MaxConn:        1,
		MaxPackageSize: 512,
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
	}
	//GlobalObject.Reload()
}