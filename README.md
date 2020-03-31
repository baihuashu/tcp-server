# 功能
- connection 代表每次conn的封装。 负责读写方法（分离）
- request 代表的是 conn和数据msg的封装
- handle 处理的是request 为了多路由，后面升级成工作池处理。
- connManger 使用map管理所有的conn 
- 连接进入后和断开连接前会调用Hook方法
- 解决tcp黏包问题
- worker池，任务队列解决高并发
# 高并发分析
客户端来一个请求 会开三个groutine
go reader() -> go DoMsgHandle() ->go writer()

假如来10w个请求，那就有10w个reader和10w的writer但他们是阻塞的不会占用cpu资源，
但是还有10w个go DoMsgHandle() 处理客户端业务会占用cpu资源。

所以开辟一个worker池，任务队列。