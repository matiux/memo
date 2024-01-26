package command

import (
	"context"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestBus_Handle(t *testing.T) {
	bus := &Bus{}

	prev := bus.Handle("1", intHandler(1))
	assert.Nil(t, prev)
	assert.Equal(t, intHandler(1), bus.handlers[reflect.TypeOf("1")])

	prev2 := bus.Handle("1", intHandler(2))
	assert.Equal(t, intHandler(1), prev2)
	assert.Equal(t, intHandler(2), bus.handlers[reflect.TypeOf("1")])
}

func TestBus_ExecuteContext_errorsWithCancelledContext(t *testing.T) {
	bus := &Bus{}

	command := intCommand(1)

	bus.Handle(command, HandlerFunc(func(ctx context.Context, cmd Command) (interface{}, error) {
		return "something we wont get later", nil
	}))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10))
	defer cancel()

	result, err := bus.ExecuteContext(ctx, command)

	assert.Nil(t, result)
	assert.Equal(t, err, context.DeadlineExceeded)
}

type intCommand int

type intHandler int

func (ih intHandler) Handle(ctx context.Context, command Command) (interface{}, error) {
	return int(ih), nil
}
