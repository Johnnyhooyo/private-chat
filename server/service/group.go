package service

import (
	"fmt"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/pkg/log"
)

// Group 群组聊天
type Group struct {
	name          string
	userPool      map[string]*common.UserInfo
	userWriteBack map[string]func(data any) error
}

// AddGroup 加入群聊
func (g *Group) AddGroup(user *common.UserInfo, receiveFunc func(data any) error) {
	for _, f := range g.userWriteBack {
		msg := common.Message{From: common.System, Body: fmt.Sprintf("用户:%s 加入群聊", user.Name)}
		_ = f(msg)
	}
	g.userPool[user.Name] = user
	g.userWriteBack[user.Name] = receiveFunc
}

// QuitGroup 退出群聊
func (g *Group) QuitGroup(user *common.UserInfo) {
	delete(g.userPool, user.Name)
	delete(g.userWriteBack, user.Name)
	for _, f := range g.userWriteBack {
		msg := common.Message{From: common.System, Body: fmt.Sprintf("用户:%s 退出群聊", user.Name)}
		_ = f(msg)
	}
}

// ReceiveMsg 群组消息转发
func (g *Group) ReceiveMsg(msg *common.Message) {
	for user, f := range g.userWriteBack {
		err := f(msg)
		if err != nil {
			log.Errorf("user:%s, error receive msg:%s", user, msg.Body)
			continue
		}
	}
}
