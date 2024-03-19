package common

import (
	"fmt"
	"github.com/johnnhooyo/private-chat/common/route"
)

type UserInfo struct {
	Ip   string `json:"ip,omitempty"`
	Name string `json:"name,omitempty"`
}

type Message struct {
	Route      route.Type `json:"route,omitempty"`
	From       *UserInfo  `json:"from,omitempty"`
	To         *UserInfo  `json:"to,omitempty"`
	Body       any        `json:"body,omitempty"` // 如果有图片之类的 用${attach_1} 格式填充，然后去attachFile中按index读取
	AttachFile [][]byte   `json:"attach_file,omitempty"`
}

func (m *Message) ReceiveMsg() {
	fmt.Printf("from:%s:%s \n", m.From.Name, m.Body)
}

func (m *Message) SendMsg() {
	fmt.Printf("to:%s:%s \n", m.To.Name, m.Body)

}

type LogResp struct {
	Logged bool   `json:"logged,omitempty"`
	ErrMsg string `json:"err_msg,omitempty"`
}

var RespMap map[route.Type]any

func init() {
	RespMap = map[route.Type]any{
		route.LogIn:    &LogResp{},
		route.LogOut:   &LogResp{},
		route.UserList: ss,
	}
}

var ss []*UserInfo
