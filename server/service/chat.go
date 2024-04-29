package service

import (
	"fmt"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/common/route"
	"github.com/johnnhooyo/private-chat/core"
	"github.com/johnnhooyo/private-chat/pkg/log"
)

func NewChatHandler() core.Handler {
	return &ChatHandler{}
}

type ChatHandler struct {
}

func (l *ChatHandler) Handle(ctx *chat.Context, req any) error {
	msg, ok := req.(*common.Message)
	if !ok {
		log.Errorf("unkonw chat message")
		return nil
	}
	toUser, ok := getGround().userPool[msg.To.Name]
	if !ok {
		if group, ok := getGround().gruops[msg.To.Name]; ok {
			group.ReceiveMsg(msg)
			return nil
		}
		log.Errorf("user %s is logout, send fail", msg.To)
		err := ctx.Write(common.Message{Route: route.Chat, From: &common.UserInfo{Name: "System"}, Body: fmt.Sprintf("user %s is logout, send fail", msg.To.Name)})
		return err
	}
	sendFunc, ok := getGround().userWriteBack[toUser.Name]
	if !ok {
		log.Errorf("user %s is logout, send fail", msg.To)
		err := ctx.Write(common.Message{Route: route.Chat, From: &common.UserInfo{Name: "System"}, Body: fmt.Sprintf("user %s is logout, send fail", msg.To.Name)})
		return err
	}
	msg.Route = route.Chat
	return sendFunc(msg)
}

func (l *ChatHandler) GetReq() core.NewRequest {
	return func() interface{} {
		return &common.Message{}
	}
}
