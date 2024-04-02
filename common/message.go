package common

import (
	"fmt"
	"github.com/johnnhooyo/private-chat/common/route"
	"os"
)

type UserInfo struct {
	Ip   string `json:"ip,omitempty"`
	Name string `json:"name,omitempty"`
}

type Message struct {
	Route      route.Type        `json:"route,omitempty"`
	From       *UserInfo         `json:"from,omitempty"`
	To         *UserInfo         `json:"to,omitempty"`
	Body       any               `json:"body,omitempty"` // 如果有图片之类的 用[pic:(.*?)] 格式填充，然后去attachFile中按index读取
	AttachFile map[string][]byte `json:"attach_file,omitempty"`
}

func (m *Message) ReceiveMsg() {
	fmt.Printf("from:%s:%s \n", m.From.Name, m.Body)
	for fileName, bytes := range m.AttachFile {
		if _, err := os.Stat(fmt.Sprintf("./%s/", m.From.Name)); os.IsNotExist(err) {
			err = os.MkdirAll(fmt.Sprintf("./%s/", m.From.Name), 0755)
			if err != nil {
				fmt.Println("无法创建目标文件夹：", err)
				return
			}
		}
		err := os.WriteFile(fmt.Sprintf("./%s/", m.From.Name)+fileName, bytes, 0644)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
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

type StringMsg struct {
	Data string
}
