package core

import (
	"encoding/binary"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/pkg/log"
	"reflect"
)

type Dispatcher interface {
	Dispatch(ctx *chat.Context, data []byte)
	Register(path string, handler Handler)
}

func NewDefaultDispatcher() Dispatcher {
	return &defaultDispatcher{
		router: make(map[string]Handler),
	}
}

type defaultDispatcher struct {
	router map[string]Handler
}

func (d *defaultDispatcher) Dispatch(ctx *chat.Context, data []byte) {
	routeSize := binary.BigEndian.Uint32(data[:4])
	route := string(data[4 : 4+routeSize])
	if handler, ok := d.router[route]; ok {
		req := handler.GetReq()()
		if req != nil {
			if err := common.InUseCodec.Unmarshal(data[4+routeSize:], req); err != nil {
				log.Errorf("err binding req data to struct %s, err is %s", reflect.TypeOf(req), err.Error())
				return
			}
		}

		err := handler.Handle(ctx, req)
		if err != nil {
			log.Errorf("handler err is %s", err.Error())
			return
		}
	} else {
		log.Warnf("unsupported route %s", route)
	}
}

func (d *defaultDispatcher) Register(path string, handler Handler) {
	d.router[path] = handler
}
