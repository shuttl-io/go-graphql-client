package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	assert := assert.New(t)
	c := NewClient(nil)
	assert.NotNil(c)
	assert.Implements((*Client)(nil), new(client))
}

func TestNewRequest(t *testing.T) {
	assert := assert.New(t)
	c := NewClient(nil)
	assert.NotNil(c.NewRequest())
}

type mockTransport struct{}

func (m *mockTransport) Transport(req Request) (Response, error) {
	return Response{}, nil
}

func TestTransport(t *testing.T) {
	assert := assert.New(t)

	m := &mockTransport{}
	c := NewClient(m).(*client)
	r := c.NewRequest().(*request)
	assert.Equal(c.transport, r.transport)

	c = NewClient(nil).(*client)
	assert.Nil(c.transport)
	c.SetTransport(m)
	r = c.NewRequest().(*request)
	assert.Equal(c.transport, r.transport)

}
