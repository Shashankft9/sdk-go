package test_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cloudevents/sdk-go/pkg/binding"
	. "github.com/cloudevents/sdk-go/pkg/binding/test"
	. "github.com/cloudevents/sdk-go/pkg/protocol/test"
)

func TestEvent(t *testing.T) {
	assert := assert.New(t)

	e := FullEvent()
	assert.Equal("1.0", e.SpecVersion())
	assert.Equal("com.example.FullEvent", e.Type())
	var s string
	err := e.DataAs(&s)
	assert.NoError(err)
	assert.Equal("hello", s)

	e = MinEvent()
	assert.Equal("1.0", e.SpecVersion())
	assert.Equal("com.example.MinEvent", e.Type())
	assert.Nil(e.Data())
	assert.Empty(e.DataContentType())
}

type dummySR chan binding.Message

func (d dummySR) Send(ctx context.Context, m binding.Message) (err error) { d <- m; return nil }
func (d dummySR) Receive(ctx context.Context) (binding.Message, error)    { return <-d, nil }

func TestSendReceive(t *testing.T) {
	sr := make(dummySR)
	allIn := []binding.Message{}
	for _, e := range Events() {
		allIn = append(allIn, binding.ToMessage(&e))
	}

	var allOut []binding.Message
	EachMessage(t, allIn, func(t *testing.T, in binding.Message) {
		SendReceive(t, context.Background(), in, sr, sr, func(out binding.Message) {
			assert.Equal(t, in, out)
			allOut = append(allOut, out)
		})
	})
	assert.Equal(t, allIn, allOut)
}