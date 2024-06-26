package main

import (
	"fmt"
	"github.com/johnnhooyo/private-chat/client/core"
	"github.com/johnnhooyo/private-chat/common/chat"
	log2 "log"
)

// 使用命令行输入来进行交互
func main() {
	startup := `welcome to private-chat, u can say whatever u want~`
	fmt.Println(startup)

	for {
		client := core.NewClient()
		// 链接到服务端 进行通信
		imClient := core.NewImClient("127.0.0.1:8002", client)
		ctx := chat.Background()
		err := imClient.Start(ctx)
		if err != nil {
			log2.Fatal("can not connect to server, err is " + err.Error())
		}

		go imClient.ReadLoop(ctx)
		// 开始交互
		ct := client.Run(ctx)
		if !ct {
			break
		}
		ctx.Cancel(fmt.Errorf("client closed"))
		var conti string
		fmt.Print("是否重新连接（y/n）: ")
		_, _ = fmt.Scanln(&conti)
		if conti != "y" {
			break
		}
	}

	fmt.Println("bye bye~~~")

}
