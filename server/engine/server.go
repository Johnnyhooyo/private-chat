// Package engine 核心引擎，负责处理链接的处理和数据的处理
package engine

import "context"

// Server 主引擎
type Server struct {
}

// Start 启动服务
func (s *Server) Start(ctx context.Context) {
	l := &listener{}
	l.listen(ctx)
}
