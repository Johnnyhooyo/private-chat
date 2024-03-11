package engine

type Options struct {
	bossNum   int
	workerNum int
}

// Option is a function that will set up option.
type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

// WithBossNum 设置主Acceptor数量
func WithBossNum(num int) Option {
	return func(opts *Options) {
		opts.bossNum = num
	}
}

// WithWorkerNum 设置从Acceptor数量
func WithWorkerNum(num int) Option {
	return func(opts *Options) {
		opts.workerNum = num
	}
}
