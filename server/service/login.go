package service

import (
	"fmt"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/core"
	"github.com/johnnhooyo/private-chat/pkg/log"
)

func NewLoginHandler() core.Handler {
	return &LoginHandler{}
}

type LoginHandler struct {
}

func (l *LoginHandler) Handle(ctx *chat.Context, req any) error {
	info, ok := req.(*common.UserInfo)
	resp := &common.LoginResp{Logged: true}
	var err error
	if !ok {
		log.Errorf("error revice login event ,req can't read")
		err = fmt.Errorf("error revice login event ,req can't read")
	} else {
		if loggedUserInfo, ok := getGround().userPool[info.Name]; ok && loggedUserInfo.Ip != info.Ip {
			err = fmt.Errorf("repeated name, please try new name")
		} else {
			getGround().userPool[info.Name] = info
		}
	}
	if err != nil {
		resp.Logged = false
		resp.ErrMsg = err.Error()
	}
	if err := ctx.Write(common.Message{Route: "login", Body: resp}); err != nil {
		log.Errorf("err send resp in login handler. %s", err.Error())
		return nil
	}
	return nil
}

func (l *LoginHandler) GetReq() core.NewRequest {
	return func() interface{} {
		return &common.UserInfo{}
	}
}
