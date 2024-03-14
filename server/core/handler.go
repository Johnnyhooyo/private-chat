package core

import (
	"github.com/johnnhooyo/private-chat/common/chat"
)

type NewRequest func() interface{}

// Handler 抽象的业务处理逻辑
type Handler interface {
	// Handle 处理链接事件的方法
	Handle(ctx *chat.Context, req any) error
	// GetReq 获取请求参数 进行转换
	GetReq() NewRequest
}

//// NewHandler 新建一个处理逻辑
//type NewHandler func() Handler

//// handlerGroup 注册的处理事件
//type handlerGroup struct {
//	router map[string]Handler
//}
//
//// Register 注册业务逻辑
//func (g *handlerGroup) Register(path string, handler Handler) {
//	g.router[path] = handler
//}
