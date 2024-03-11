package engine

import (
	"net"

	"github.com/johnnhooyo/private-chat/pkg/chat"
	"github.com/johnnhooyo/private-chat/pkg/log"
)

// acceptorGroup 连接处理器组
type acceptorGroup struct {
	group []*acceptor
}

func (a *acceptorGroup) listen(ctx *chat.Context, lis net.Listener) {

}

// 连接处理器
type acceptor struct {
}

func listen(ctx *chat.Context, lr net.Listener) {
	for {
		conn, err := lr.Accept()
		if err != nil {
			log.Errorf("accept connection err %s", err)
			continue
		}
		go func() {
			h := &handler{}
			h.handle(ctx, conn)
		}()
	}
}
