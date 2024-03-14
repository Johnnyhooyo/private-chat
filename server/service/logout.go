package service

import (
	"fmt"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/core"
	"github.com/johnnhooyo/private-chat/pkg/log"
)

func NewLogoutHandler() core.Handler {
	return &LogoutHandler{}
}

type LogoutHandler struct {
}

func (l *LogoutHandler) Handle(_ *chat.Context, req any) error {
	info, ok := req.(*common.UserInfo)
	if !ok {
		log.Errorf("error revice logout event ,req can't read")
		return fmt.Errorf("unknow req type")
	}
	delete(getGround().userPool, info.Name)
	return nil
}

func (l *LogoutHandler) GetReq() core.NewRequest {
	return func() interface{} {
		return &common.UserInfo{}
	}
}
