package core

import (
	"encoding/binary"
	"github.com/johnnhooyo/private-chat/client/log"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"net"
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
			log.Errorf("failed read data from server error:%s", err.Error())
			continue
		}
		routeSize := binary.BigEndian.Uint32(bytes[:4])
		route := string(bytes[4 : 4+routeSize])

		msg := &common.Message{Route: route}
		err = common.InUseCodec.Unmarshal(bytes[4+routeSize:], msg)
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
