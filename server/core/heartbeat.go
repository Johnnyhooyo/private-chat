package core

import (
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/common/route"
	"github.com/johnnhooyo/private-chat/pkg/log"
	"github.com/panjf2000/gnet/v2"
	"sync"
	"time"
)

var once sync.Once

// GetHeartbeat 新建一个心跳处理
func GetHeartbeat() *heartbeat {
	once.Do(func() {
		hb = &heartbeat{
			rwLock: sync.RWMutex{},
			cm:     make(map[int]time.Time),
		}
	})
	return hb
}

var hb *heartbeat

type heartbeat struct {
	rwLock sync.RWMutex
	cm     map[int]time.Time
}

// CheckConn 添加一个需要监控的链接
func (h *heartbeat) CheckConn(ctx *chat.Context, conn gnet.Conn) {
	h.rwLock.Lock()
	defer h.rwLock.Unlock()

	if lastBeat, exists := h.cm[conn.Fd()]; exists {
		if lastBeat.Add(3 * time.Second).Before(time.Now()) {
			err := conn.Close()
			if err != nil {
				log.Errorf("conn close error:%s", err.Error())
			} else {
				delete(h.cm, conn.Fd())
			}
			return
		}
	} else {
		h.cm[conn.Fd()] = time.Now()
	}
	_ = ctx.Write(common.Message{
		Route: route.Heartbeat,
		Body:  common.StringMsg{Data: "ping"},
	})
}

// Reset 重置心跳超时
func (h *heartbeat) Reset(fd int) {
	h.rwLock.Lock()
	defer h.rwLock.Unlock()

	h.cm[fd] = time.Now()
}
