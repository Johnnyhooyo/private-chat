package main

import (
	"fmt"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/core"
	"github.com/johnnhooyo/private-chat/service"
	"github.com/panjf2000/gnet/v2"
)

func main() {
	//server := engine.NewServer()
	//server.Start(chat.Background())
	server := core.NewImServer(common.GenerateAESKey("ucanmodifythisword"))
	server.Register("login", service.NewLoginHandler())
	server.Register("logout", service.NewLogoutHandler())
	server.Register("userlist", service.NewUserListHandler())
	server.Register("chat", service.NewChatHandler())
	server.Register("heartbeat", service.NewHeartbeatHandler())

	err := gnet.Run(server, "tcp://0.0.0.0:8002", gnet.WithTicker(true))
	if err != nil {
		fmt.Printf("start gnet error :%s", err.Error())
		return
	}
}
