# private-chat
Private chat tool
计划实现一个私有的聊天工具，服务端进行数据转发，客户端使用命令行进行聊天。

client -- tcp  --- server
server {
    engine 核心引擎
    codec 编解码
    package 数据包处理
    router 请求路由
    service 业务处理
}