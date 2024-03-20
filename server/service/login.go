package service

import (
	"fmt"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/common/route"
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
	resp := &common.LogResp{Logged: true}
	var err error
	if !ok {
		log.Errorf("error revice login event ,req can't read")
		err = fmt.Errorf("error revice login event ,req can't read")
	} else {
		if loggedUserInfo, ok := getGround().userPool[info.Name]; ok && loggedUserInfo.Ip != info.Ip {
			err = fmt.Errorf("repeated name, please try new name")
		} else {
			getGround().userPool[info.Name] = info
			getGround().userWriteBack[info.Name] = ctx.Write
			log.Infof("user %s loggedin from ip:%s", info.Name, info.Ip)
		}
	}
	if err != nil {
		resp.Logged = false
		resp.ErrMsg = err.Error()
	}
	if err := ctx.Write(common.Message{Route: route.LogIn, Body: resp}); err != nil {
		log.Errorf("err send resp in login handler. %s", err.Error())
		return nil
	}
	err = ctx.Broadcast(common.Message{
		Route: route.Broadcast,
		From:  info,
		Body:  resp,
	})
	if err != nil {
		log.Errorf("err send broadcast in login handler. %s", err.Error())
		return err
	}
	return nil
}

func (l *LoginHandler) GetReq() core.NewRequest {
	return func() interface{} {
		return &common.UserInfo{}
	}
}
