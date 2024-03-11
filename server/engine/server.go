// Package engine 核心引擎，负责处理链接的处理和数据的处理
package engine

import "github.com/johnnhooyo/private-chat/pkg/chat"

// NewServer new server
func NewServer(opts ...Option) *Server {
	Options := loadOptions(opts...)
	bossAcceptor := &acceptorGroup{
		group: make([]*acceptor, Options.bossNum),
	}

	listener := &listener{
		addr:      "",
		bossGroup: bossAcceptor,
	}

	return &Server{
		lr: listener,
	}
}

// Server 主引擎
type Server struct {
	lr *listener
}

// Start 启动服务
func (s *Server) Start(ctx *chat.Context) {
	l := &listener{
		addr: ":8080",
	}
	l.listen(ctx)
}
