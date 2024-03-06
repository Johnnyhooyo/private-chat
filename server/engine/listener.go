package engine

import (
	"context"
	"log"
	"net"

	logger "github.com/johnnhooyo/private-chat/pkg/log"
)

type listener struct {
	addr string
}

func (l *listener) listen(ctx context.Context) error {
	lr, err := net.Listen("tcp", l.addr)
	if err != nil {
		log.Fatal("listen to port error", l.addr)
	}
	for {
		conn, err := lr.Accept()
		if err != nil {
			logger.Errorf("accept connection err %s", err)
			continue
		}
		go func() {
			h := &handler{}
			h.handle(ctx, conn)
		}()
	}
}
