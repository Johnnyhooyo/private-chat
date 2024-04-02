# private-chat
Private chat tool
计划实现一个私有的聊天工具，服务端进行数据转发，客户端使用命令行进行聊天。

```
client <--tcp--> server <--tcp--> client
server {
    engine 核心引擎
    codec 编解码
    package 数据包处理
    router 请求路由
    service 业务处理
}
client {
    input 输入模块
    listen 监听模块
    codec 编解码
    package 数据包处理
}
```

如果你是前端大佬，可以做个页面出来。


# todo:
- 图片、视频文件的压缩
- 代码优化、整理
