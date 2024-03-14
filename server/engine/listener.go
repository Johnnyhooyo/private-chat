package engine

import (
	"github.com/johnnhooyo/private-chat/common/chat"
	"log"
	"net"
)

type listener struct {
	addr      string
	bossGroup *acceptorGroup
}

func (l *listener) listen(ctx *chat.Context) {
	lr, err := net.Listen("tcp", l.addr)
	if err != nil {
		log.Fatal("listen to port error", l.addr)
	}
	l.bossGroup.listen(ctx, lr)
}
