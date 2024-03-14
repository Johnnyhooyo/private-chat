module github.com/johnnhooyo/private-chat/client

go 1.20

require (
	github.com/johnnhooyo/private-chat/common v0.0.0
	go.uber.org/zap v1.27.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

replace github.com/johnnhooyo/private-chat/common => ../common

require (
	github.com/stretchr/testify v1.8.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)
