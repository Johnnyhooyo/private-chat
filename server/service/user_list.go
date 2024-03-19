package service

import (
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/common/route"
	"github.com/johnnhooyo/private-chat/core"
)

func NewUserListHandler() core.Handler {
	return &UserListHandler{}
}

type UserListHandler struct {
}

// Handle 返回用户列表 暂时没有排除自己
func (l *UserListHandler) Handle(ctx *chat.Context, req any) error {
	var userList []*common.UserInfo
	for _, info := range getGround().userPool {
		userList = append(userList, info)
	}
	err := ctx.Write(common.Message{Route: route.UserList, Body: userList})
	if err != nil {
		return err
	}
	return nil
}

func (l *UserListHandler) GetReq() core.NewRequest {
	return func() interface{} {
		return nil
	}
}
