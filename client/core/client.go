package core

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/johnnhooyo/private-chat/common"
	"github.com/johnnhooyo/private-chat/common/chat"
	"github.com/johnnhooyo/private-chat/common/route"
	"strings"
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
	}
}

func (c *Client) HandleMsg(msg *common.Message) error {
	switch msg.Route {
	case route.Chat:
		msg.ReceiveMsg()
	case route.LogIn:
		if resp, ok := msg.Body.(common.LoginResp); ok {
			if resp.Logged {
				c.logged = true
				c.request[route.LogIn] <- true
			}
		}
	case route.LogOut:

	}
	return nil
}

func (c *Client) Run(ctx *chat.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ctx canceled, exit now.")
			return
		default:

		}

		if !c.logged {
			c.Login()
			fmt.Println("you can input 'help' to get more info.")
		}

		var input string
		_, _ = fmt.Scanln(&input)
		c.cmd(ctx, input)
	}
}

func (c *Client) cmd(ctx *chat.Context, msg string) {
	messages := strings.Split(msg, "_")
	cmd := messages[0]
	var user string
	var body string
	if len(messages) >= 1 {
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
	case "user_list":
		c.GetUserList()
	case "to":
		c.SendMsg(user, body)
	case "reject":
		if err := c.Request(route.Reject, &common.UserInfo{Name: user}); err != nil {
			fmt.Printf("处理错误，请重试。")
		}
	case "recover":
		if err := c.Request(route.Recover, &common.UserInfo{Name: user}); err != nil {
			fmt.Printf("处理错误，请重试。")
		}
	default:
		fmt.Printf("暂不支持的指令，如果需要可以输入“help”获取帮助")
	}
}

func (c *Client) Request(routeStr route.Type, message any) error {
	routeBytes := []byte(routeStr)
	var msg []byte
	binary.BigEndian.AppendUint32(msg, uint32(len(routeBytes)))
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

func (c *Client) Login() {
	fmt.Print("请输入你的名字: ")
	_, _ = fmt.Scanln(&c.name)

	user := &common.UserInfo{
		Name: c.name,
	}
	err := c.Request(route.LogIn, user)
	if err != nil {
		fmt.Printf("处理错误，请重试。")
		return
	}
	select {
	case <-c.request[route.LogIn]:

	}
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
	select {
	case <-c.request[route.UserList]:

	}
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
