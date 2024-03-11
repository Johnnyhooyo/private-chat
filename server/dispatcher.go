package main

type Dispatcher interface {
	Dispatch(data []byte)
}

func NewDefaultDispatcher() Dispatcher {
	return &defaultDispatcher{}
}

type defaultDispatcher struct {
}

func (d *defaultDispatcher) Dispatch(data []byte) {
	//TODO implement me
	panic("implement me")
}
