package service

import (
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/common/route"
	"github.com/johnnhooyo/private-chat/core"
)

// NewHeartbeatHandler 获取心跳处理
func NewHeartbeatHandler() core.Handler {
	return &HeartbeatHandler{}
}

// HeartbeatHandler 心跳处理 也可以在gnet中开启 so_keepalive 选项 有gnet来处理心跳
type HeartbeatHandler struct {
}

// Handle 心跳监测
func (l *HeartbeatHandler) Handle(ctx *chat.Context, _ any) error {
	hb := core.GetHeartbeat()
	v, ok := ctx.Get(string(route.Heartbeat))
	if ok {
		hb.Reset(v.(int))
	}
	return nil
}

func (l *HeartbeatHandler) GetReq() core.NewRequest {
	return func() interface{} {
		return &common.StringMsg{}
	}
}
