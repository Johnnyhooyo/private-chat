package main

import (
	"context"

	"github.com/johnnhooyo/private-chat/engine"
)

func main() {
	server := &engine.Server{}
	server.Start(context.Background())
}
