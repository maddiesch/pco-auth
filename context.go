package auth

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// Context manages the state of an auth request.
type authContext struct {
	ctx  context.Context
	mu   sync.Mutex
	done chan struct{}
	err  error
}

// Failable is a context type that can fail with an error.
type Failable interface {
	Fail(err error)
}

func withContext(ctx context.Context) *authContext {
	new := &authContext{
		ctx:  ctx,
		done: make(chan struct{}),
	}

	go func(new *authContext, ctx context.Context) {
		if ctx.Done() == nil { // The sub-context can never be cancelled
			return
		}

		<-ctx.Done()

		new.mu.Lock()
		new.err = ctx.Err()
		new.done <- struct{}{}
		new.mu.Unlock()
	}(new, ctx)

	return new
}

// Fail sets the context error and exits the goroutine
func (c *authContext) Fail(err error) {
	c.mu.Lock()
	c.err = err
	c.done <- struct{}{}
	c.mu.Unlock()
	runtime.Goexit()
}

// Deadline returns the completion deadline. (None for this context)
func (c *authContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Done is called if the context receives a complete signal.
func (c *authContext) Done() <-chan struct{} {
	c.mu.Lock()
	d := c.done
	c.mu.Unlock()

	return d
}

// Err returns the error from the context
func (c *authContext) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

// Value is the context's value (nil for this context)
func (c *authContext) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// String returns the string representation
func (c *authContext) String() string {
	return "Context"
}

func _fail(ctx context.Context, err error) {
	if ctx, ok := ctx.(Failable); ok {
		ctx.Fail(err)
	} else {
		panic(err)
	}
}
