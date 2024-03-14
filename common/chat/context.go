// Package chat 类似于gin的context实现
package chat

import (
	"sync"
	"time"
)

// Background 初始化一个空的Context
func Background() *Context {
	return &Context{
		mu:  sync.RWMutex{},
		rmu: sync.Mutex{},
	}
}

// Context 串联空间
type Context struct {
	// This mutex protects Keys map and done chan.
	mu sync.RWMutex
	// Keys is a key/value pair exclusively for the context of each request.
	Keys map[string]any

	// This mutex protects Conn write and broadcast with goroutine-safe.
	rmu sync.Mutex
	//// Conn to reply msg or push msg.
	//Conn gnet.Conn
	write func(data any) error

	// broadcast some system message.
	broadcast func(data any) error

	// for cancel
	done chan struct{}
	err  error
}

// Cancel cancel context
func (c *Context) Cancel(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.done == nil {
		c.done = make(chan struct{})
		c.err = err
		close(c.done)
	}
}

//// BindConn bind a gnet.Conn to this context.
//func (c *Context) BindConn(conn gnet.Conn) {
//	c.Conn = conn
//}

// BindWriteFunc 数据写回到通道
func (c *Context) BindWriteFunc(f func(data any) error) {
	c.write = f
}

// Write 数据写回到通道
func (c *Context) Write(data interface{}) error {
	c.rmu.Lock()
	defer c.rmu.Unlock()
	return c.write(data)
}

// BindBroadcastFunc 绑定广播通道
func (c *Context) BindBroadcastFunc(f func(data any) error) {
	c.broadcast = f
}

func (c *Context) Broadcast(data any) error {
	c.rmu.Lock()
	defer c.rmu.Unlock()
	return c.broadcast(data)
}

// Set is used to store a new key/value pair exclusively for this context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]any)
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exist it returns (nil, false)
func (c *Context) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

///////////////////////////
// impl context.Context	///
///////////////////////////

// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	panic("not implemented") // TODO: Implement
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value.
// The close of the Done channel may happen asynchronously,
// after the cancel function returns.
//
// WithCancel arranges for Done to be closed when cancel is called;
// WithDeadline arranges for Done to be closed when the deadline
// expires; WithTimeout arranges for Done to be closed when the timeout
// elapses.
//
// Done is provided for use in select statements:
//
//	// Stream generates values with DoSomething and sends them to out
//	// until DoSomething returns an error or ctx.Done is closed.
//	func Stream(ctx context.Context, out chan<- Value) error {
//		for {
//			v, err := DoSomething(ctx)
//			if err != nil {
//				return err
//			}
//			select {
//			case <-ctx.Done():
//				return ctx.Err()
//			case out <- v:
//			}
//		}
//	}
//
// See https://blog.golang.org/pipelines for more examples of how to use
// a Done channel for cancellation.
func (c *Context) Done() <-chan struct{} {
	return c.done
}

// If Done is not yet closed, Err returns nil.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the context was canceled
// or DeadlineExceeded if the context's deadline passed.
// After Err returns a non-nil error, successive calls to Err return the same error.
func (c *Context) Err() error {
	panic("not implemented") // TODO: Implement
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
//
// Use context values only for request-scoped data that transits
// processes and API boundaries, not for passing optional parameters to
// functions.
//
// A key identifies a specific value in a Context. Functions that wish
// to store values in Context typically allocate a key in a global
// variable then use that key as the argument to context.WithValue and
// Context.Value. A key can be any type that supports equality;
// packages should define keys as an unexported type to avoid
// collisions.
//
// Packages that define a Context key should provide type-safe accessors
// for the values stored using that key:
//
//	// Package user defines a User type that's stored in Contexts.
//	package user
//
//	import "context"
//
//	// User is the type of value stored in the Contexts.
//	type User struct {...}
//
//	// key is an unexported type for keys defined in this package.
//	// This prevents collisions with keys defined in other packages.
//	type key int
//
//	// userKey is the key for user.User values in Contexts. It is
//	// unexported; clients use user.NewContext and user.FromContext
//	// instead of using this key directly.
//	var userKey key
//
//	// NewContext returns a new Context that carries value u.
//	func NewContext(ctx context.Context, u *User) context.Context {
//		return context.WithValue(ctx, userKey, u)
//	}
//
//	// FromContext returns the User value stored in ctx, if any.
//	func FromContext(ctx context.Context) (*User, bool) {
//		u, ok := ctx.Value(userKey).(*User)
//		return u, ok
//	}
func (c *Context) Value(key any) any {
	panic("not implemented") // TODO: Implement
}
