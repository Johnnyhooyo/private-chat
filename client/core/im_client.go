package core

import (
	"errors"
	"fmt"
	"github.com/johnnhooyo/private-chat/client/log"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"io"
	"net"
	"runtime/debug"
)

func NewImClient(addr string, client *Client) *ImClient {
	return &ImClient{
		addr:   addr,
		packer: common.NewDefaultPacker(),
		client: client,
	}
}

type ImClient struct {
	addr   string
	packer common.Packer
	client *Client
	c      net.Conn
}

func (i *ImClient) Start(_ *chat.Context) (err error) {
	i.c, err = net.Dial("tcp", i.addr)
	if err != nil {
		log.Errorf("connect to server error:%s", err.Error())
		return err
	}
	log.Infof("connection to server %s starts...", i.c.RemoteAddr().String())
	i.client.write = func(msg []byte) error {
		pack, err := i.packer.Pack(msg)
		if err != nil {
			return err
		}
		_, err = i.c.Write(pack)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (i *ImClient) ReadLoop(ctx *chat.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 处理 panic 错误，例如打印错误信息或进行日志记录等操作
			fmt.Println("Recovered from panic:", err)
			debug.PrintStack() // 打印出错时的调用栈信息
			i.ReadLoop(ctx)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			log.Infof("context canceled, stop read from conn.")
			_ = i.c.Close()
			return
		default:
			// continue
		}
		bytes, err := i.packer.Unpack(i.c)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Errorf("failed read data from server error:%s", err.Error())
			continue
		}

		msg := &common.Message{}
		err = common.InUseCodec.Unmarshal(bytes, msg)
		log.Debugf("receive msg %+v", msg)
		if err != nil {
			log.Errorf("failed convert data to common.Message, err is %s", err.Error())
			continue
		}
		err = i.client.HandleMsg(msg)
		if err != nil {
			log.Errorf("failed handle common.Message, err is %s", err.Error())
		}
	}
}
