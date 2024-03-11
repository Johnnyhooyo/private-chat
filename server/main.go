package main

import (
	"github.com/johnnhooyo/private-chat/engine"
	"github.com/johnnhooyo/private-chat/pkg/chat"
	"github.com/panjf2000/gnet/v2"
)

func main() {
	server := engine.NewServer()
	server.Start(chat.Background())

	gnet.Run(NewImServer(), "tcp://:8088")
}
