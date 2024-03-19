package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/johnnhooyo/private-chat/client/log"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/common/route"
	"runtime/debug"
	"strings"
	"time"
)

type Client struct {
	name         string
	logged       bool // 0-未登录 1-已登陆
	chatBoxes    map[string]Handler
	write        func(msg []byte) error
	loginVersion int64 // 登陆的事件戳 如果不匹配就抛弃消息

	request map[route.Type]chan bool // 已经发送请求 等待结果
}

func NewClient() *Client {
	return &Client{
		logged:    false,
		chatBoxes: make(map[string]Handler),
		request:   make(map[route.Type]chan bool),
	}
}

func (c *Client) HandleMsg(msg *common.Message) error {
	switch msg.Route {
	case route.Chat:
		msg.ReceiveMsg()
	case route.LogIn:
		bytes, err := common.InUseCodec.Marshal(msg.Body)
		if err != nil {
			return err
		}
		req := common.RespMap[route.LogIn]
		err = common.InUseCodec.Unmarshal(bytes, req)
		if err != nil {
			return err
		}
		if resp, ok := req.(*common.LogResp); ok {
			if resp.Logged {
				c.logged = true
				log.Debugf("send loggedin singnal into chan %+v", c.request[route.LogIn])
				c.request[route.LogIn] <- true
				fmt.Println("u r logged in.")
			}
		}
	case route.LogOut:
		// 已经下线 不用处理
	case route.UserList:
		bytes, err := common.InUseCodec.Marshal(msg.Body)
		if err != nil {
			return err
		}
		req := common.RespMap[route.UserList]
		if resp, ok := req.([]*common.UserInfo); ok {
			err = common.InUseCodec.Unmarshal(bytes, &resp)
			if err != nil {
				return err
			}
			c.request[route.UserList] <- true
			fmt.Println("当前在线的用户有：")
			for i, userInfo := range resp {
				fmt.Printf("%d: %s\n", i, userInfo.Name)
			}
			fmt.Println("-------------------")

		} else {
			log.Errorf("error resp type %+v", req)
		}

	case route.Broadcast:
		fmt.Printf("用户 %s 上线了\n", msg.From.Name)
	}
	return nil
}

func (c *Client) Run(ctx *chat.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 处理 panic 错误，例如打印错误信息或进行日志记录等操作
			fmt.Println("Recovered from panic:", err)
			debug.PrintStack() // 打印出错时的调用栈信息
			c.Run(ctx)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ctx canceled, exit now.")
			return
		default:

		}

		if !c.logged {
			if !c.Login() {
				fmt.Println("you can input 'help' to get more info.")
				continue
			}
		}

		var input string
		_, _ = fmt.Scanln(&input)
		c.cmd(ctx, input)
	}
}

func (c *Client) cmd(ctx *chat.Context, msg string) {
	if len(msg) == 0 {
		return
	}
	messages := strings.Split(msg, "_")
	cmd := messages[0]
	var user string
	var body string
	if len(messages) > 1 {
		content := strings.Split(messages[1], ":")
		user = content[0]
		if len(content) > 0 {
			body = content[1]
		}
	}
	switch cmd {
	case "help":
		PrintHelper()
	case "exit":
		c.Logout()
		if !c.logged {
			ctx.Cancel(errors.New("user exit"))
		}
	case "userlist":
		c.GetUserList()
	case "to":
		c.SendMsg(user, body)
	case "reject":
		if err := c.Request(route.Reject, &common.UserInfo{Name: user}); err != nil {
			fmt.Printf("处理错误，请重试。\n")
		}
	case "recover":
		if err := c.Request(route.Recover, &common.UserInfo{Name: user}); err != nil {
			fmt.Printf("处理错误，请重试。\n")
		}
	default:
		fmt.Printf("暂不支持的指令，如果需要可以输入“help”获取帮助\n")
	}
}

func (c *Client) Request(routeStr route.Type, message any) error {
	routeBytes := []byte(routeStr)
	var msg []byte
	msg = binary.BigEndian.AppendUint32(msg, uint32(len(routeBytes)))
	msg = append(msg, routeBytes...)
	userInfoBytes, err := common.InUseCodec.Marshal(message)
	if err != nil {
		return err
	}
	msg = append(msg, userInfoBytes...)
	err = c.write(msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Login() bool {
	fmt.Print("请输入你的名字: ")
	_, _ = fmt.Scanln(&c.name)

	user := &common.UserInfo{
		Name: c.name,
	}
	err := c.Request(route.LogIn, user)
	if err != nil {
		fmt.Printf("处理错误，请重试。")
		return false
	}
	return c.waitResp(route.LogIn)
}

func (c *Client) waitResp(types route.Type) bool {
	timeout := time.After(2 * time.Second)
	waitChannel := make(chan bool)
	c.request[types] = waitChannel
	select {
	case <-waitChannel:
		return true
	case <-timeout:
		fmt.Printf("服务器无响应，请重试：")
	}
	return false
}

func (c *Client) Logout() {
	user := &common.UserInfo{
		Name: c.name,
	}
	err := c.Request(route.LogOut, user)
	if err != nil {
		fmt.Printf("处理错误，请重试。")
		return
	}
	c.logged = false
}

func (c *Client) GetUserList() {
	err := c.Request(route.UserList, nil)
	if err != nil {
		fmt.Printf("处理错误，请重试。")
		return
	}
	c.waitResp(route.UserList)
}

func (c *Client) SendMsg(msg, user string) {
	body := &common.Message{
		To:   &common.UserInfo{Name: user},
		Body: msg,
	}
	err := c.Request(route.Chat, body)
	if err != nil {
		fmt.Printf("发送失败,请重试")
		return
	}
}
