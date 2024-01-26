package command

// https://github.com/gogolfing/cbus

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

// HandlerNotFoundError occurs when a Handler has not been registered for a Command's type.
type HandlerNotFoundError struct {
	//Command is the Command that a Handler was not found for.
	Command
}

// Error is the error implementation.
func (e *HandlerNotFoundError) Error() string {
	return fmt.Sprintf("cbus: Handler not found for Command type %T", e.Command)
}

// ExecutionPanicError occurs when a Handler or Listener panics during execution.
type ExecutionPanicError struct {
	//Panic is the value received from recover() if not nil.
	Panic interface{}
}

// Error is the error implementation.
func (e *ExecutionPanicError) Error() string {
	return fmt.Sprintf("cbus: panic while executing command %v", e.Panic)
}

//---------------------------------------------------------------------------------------------
//---------------------------------------------------------------------------------------------

// Command is an empty interface that anything can implement and allows for executing arbitrary values on a Bus.
// Therefore, a Command is any defined type that get associated with Handler.
// The specific implementation of a Command can then carry the payload for the command to execute.
type Command interface{}

// Handler defines the contract for executing a Command within a context.Context.
// The result and err return parameters will be returned from Bus.Execute*() calls which allows Command executors to
// know the results of the Command's execution.
type Handler interface {
	Handle(ctx context.Context, command Command) (result interface{}, err error)
}

// HandlerFunc is a function definition for a Handler.
type HandlerFunc func(ctx context.Context, command Command) (result interface{}, err error)

// Handle calls hf with ctx and command.
func (hf HandlerFunc) Handle(ctx context.Context, command Command) (result interface{}, err error) {
	return hf(ctx, command)
}

// Bus is the Command Bus implementation.
// A Bus contains a one to one mapping from Command types to Handlers.
// The reflect.TypeOf interface is used as keys to map from Commands to Handlers.
// All Command Handlers is called from a Bus in a newly spawned goroutine per Command execution.
// The zero value for Bus is fully functional.
// Type Bus is safe for use by multiple goroutines.
type Bus struct {
	//lock protects all other fields in Bus.
	lock     sync.RWMutex
	handlers map[reflect.Type]Handler
}

// Handle associates a Handler in b that will be called when a Command whose type equals command's type.
// Only one Handler is allowed per Command type. Any previously added Handlers
// with the same commandType will be overwritten.
// prev is the Handler previously associated with commandType if it exists.
func (b *Bus) Handle(command Command, handler Handler) (prev Handler) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.putHandler(reflect.TypeOf(command), handler)
}

func (b *Bus) putHandler(commandType reflect.Type, handler Handler) Handler {
	prev := b.handlers[commandType]
	if b.handlers == nil {
		b.handlers = map[reflect.Type]Handler{}
	}
	b.handlers[commandType] = handler
	return prev
}

// Execute is sugar for b.ExecuteContext(context.Background(), command).
func (b *Bus) Execute(command Command) (result interface{}, err error) {
	return b.ExecuteContext(context.Background(), command)
}

// ExecuteContext attempts to find a Handler for command's Type().
// If a Handler is not found, then ErrHandlerNotFound is returned immediately.
// If a Handler is found, then a new goroutine is spawned with command's Handler
func (b *Bus) ExecuteContext(ctx context.Context, command Command) (result interface{}, err error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	handler, ok := b.handlers[reflect.TypeOf(command)]
	if !ok {
		return nil, &HandlerNotFoundError{command}
	}

	return b.execute(ctx, command, handler)
}

func (b *Bus) execute(ctx context.Context, command Command, handler Handler) (interface{}, error) {
	done := make(chan *executePayload)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- &executePayload{nil, &ExecutionPanicError{r}}
			}
		}()

		result, err := handler.Handle(ctx, command)
		done <- &executePayload{result, err}
	}()

	payload := &executePayload{}
	select {
	case <-ctx.Done():
		payload.err = ctx.Err()
	case payload = <-done:
	}

	return payload.result, payload.err
}

type executePayload struct {
	result interface{}
	err    error
}
