package service

import (
	"github.com/johnnhooyo/private-chat/common"
	"sync"
)

type UserGround struct {
	userPool map[string]*common.UserInfo
}

var (
	once          sync.Once
	defaultGround *UserGround
)

// getGround 操作用户广场
func getGround() *UserGround {
	once.Do(func() {
		defaultGround = &UserGround{userPool: make(map[string]*common.UserInfo)}
	})
	return defaultGround
}
