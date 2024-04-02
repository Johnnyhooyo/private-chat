package core

import (
	"fmt"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/common/route"
	"github.com/johnnhooyo/private-chat/pkg/log"
	"github.com/panjf2000/gnet/v2"
	"time"
)

func NewImServer(key []byte) *ImServer {
	return &ImServer{
		BuiltinEventEngine: BuiltinEventEngine{},
		eng:                gnet.Engine{},
		connMap:            make(map[int]gnet.Conn),
		packer:             common.NewDefaultPacker(key),
		dispatcher:         NewDefaultDispatcher(),
		hb:                 GetHeartbeat(),
	}
}

type ImServer struct {
	BuiltinEventEngine
	eng        gnet.Engine
	connMap    map[int]gnet.Conn
	packer     common.Packer
	dispatcher Dispatcher
	hb         *heartbeat
}

func (i *ImServer) Register(path string, handler Handler) {
	i.dispatcher.Register(path, handler)
}

func (i *ImServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	i.eng = eng
	return
}

func (i *ImServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	i.connMap[c.Fd()] = c
	log.Debugf("new connection opened from %s, local addr:%s", c.RemoteAddr(), c.LocalAddr())
	return
}

func (i *ImServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if _, ok := i.connMap[c.Fd()]; ok {
		delete(i.connMap, c.Fd())
	}
	log.Debugf("connection closed now <%s> err is %+v", c.RemoteAddr(), err)
	return
}

// OnTraffic 从conn中读取数据 并进行处理
func (i *ImServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	defer func() {
		if err := recover(); err != nil {
			// 处理 panic 错误，例如打印错误信息或进行日志记录等操作
			fmt.Println("Recovered from panic:", err)
		}
	}()

	dataPack, err := i.packer.Unpack(c)
	if err != nil {
		log.Errorf("onTraffic unpack message is error: %s", err.Error())
		return
	}
	ctx := chat.Background()
	ctx.Set(string(route.Heartbeat), c.Fd())
	ctx.BindWriteFunc(writeFunc(c, i))
	// 添加广播 通知除了自己的其他在线用户
	ctx.BindBroadcastFunc(func(data any) error {
		if bytes, err := common.InUseCodec.Marshal(data); err != nil {
			return err
		} else {
			if packageData, err := i.packer.Pack(bytes); err != nil {
				return err
			} else {
				for _, conn := range i.connMap {
					if conn.Fd() == c.Fd() {
						continue
					}

					n, err := conn.Write(packageData)
					if err != nil {
						return err
					}
					if n != len(packageData) {
						return fmt.Errorf("data not send enough, send len %d", n)
					}
				}
			}
		}
		return nil
	})
	i.dispatcher.Dispatch(ctx, dataPack)
	return
}

func writeFunc(c gnet.Conn, i *ImServer) func(data any) error {
	return func(data any) error {
		log.Debugf("send back to client %s data %+v", c.RemoteAddr(), data)
		if bytes, err := common.InUseCodec.Marshal(data); err != nil {
			return err
		} else {
			if packageData, err := i.packer.Pack(bytes); err != nil {
				return err
			} else {
				n, err := c.Write(packageData)
				if err != nil {
					return err
				}
				if n != len(packageData) {
					return fmt.Errorf("data not send enough, send len %d", n)
				}
			}
		}
		return nil
	}
}

func (i *ImServer) OnTick() (delay time.Duration, action gnet.Action) {
	for _, conn := range i.connMap {
		ctx := chat.Background()
		ctx.BindWriteFunc(writeFunc(conn, i))
		i.hb.CheckConn(ctx, conn)
	}
	return time.Second, gnet.None
}

// BuiltinEventEngine is a built-in implementation of EventHandler which sets up each method with a default implementation,
// you can compose it with your own implementation of EventHandler when you don't want to implement all methods
// in EventHandler.
type BuiltinEventEngine struct{}

// OnBoot fires when the engine is ready for accepting connections.
// The parameter engine has information and various utilities.
func (*BuiltinEventEngine) OnBoot(_ gnet.Engine) (action gnet.Action) {
	return
}

// OnShutdown fires when the engine is being shut down, it is called right after
// all event-loops and connections are closed.
func (*BuiltinEventEngine) OnShutdown(_ gnet.Engine) {
}

// OnOpen fires when a new connection has been opened.
// The parameter out is the return value which is going to be sent back to the peer.
func (*BuiltinEventEngine) OnOpen(_ gnet.Conn) (out []byte, action gnet.Action) {
	return
}

// OnClose fires when a connection has been closed.
// The parameter err is the last known connection error.
func (*BuiltinEventEngine) OnClose(_ gnet.Conn, _ error) (action gnet.Action) {
	return
}

// OnTraffic fires when a local socket receives data from the peer.
func (*BuiltinEventEngine) OnTraffic(_ gnet.Conn) (action gnet.Action) {
	return
}

// OnTick fires immediately after the engine starts and will fire again
// following the duration specified by the delay return value.
func (*BuiltinEventEngine) OnTick() (delay time.Duration, action gnet.Action) {
	return
}
