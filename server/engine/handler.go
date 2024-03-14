package engine

import (
	"github.com/johnnhooyo/private-chat/common/chat"
	"net"

	"github.com/johnnhooyo/private-chat/pkg/log"
)

type handler struct {
}

func (h *handler) handle(ctx *chat.Context, conn net.Conn) error {
	for {
		buffer := [1024]byte{}
		len, err := conn.Read(buffer[:])
		if err != nil {
			log.Errorf("read conn %s err %s", conn.RemoteAddr(), err)
		}
		msg := string(buffer[:len])
		log.Debugf("msg is %s", msg)
	}
}
