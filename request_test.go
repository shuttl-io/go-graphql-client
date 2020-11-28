package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanBuildASimpleQuery(t *testing.T) {
	assert := assert.New(t)
	req := newReq()
	req.Query(&struct {
		Example struct {
			Message string `json:"message"`
			Boolean bool   `json:"boolean"`
			Number  int    `json:"number"`
		} `json:"example" gql_params:"name:Int"`
	}{})
	reqStr := req.GetQuery()
	assert.Equal(`query{
    example{
        message
        boolean
        number
    }
}
`, reqStr)
}

func TestCanBuildASimpleMutation(t *testing.T) {
	assert := assert.New(t)
	req := newReq()
	req.Mutation(&struct {
		Example struct {
			Message string `json:"message"`
			Boolean bool   `json:"boolean"`
			Number  int    `json:"number"`
		} `json:"example" gql_params:"name:Int"`
	}{})
	reqStr := req.GetQuery()
	assert.Equal(`mutation{
    example{
        message
        boolean
        number
    }
}
`, reqStr)
}

func TestCanBuildASimpleQueryWithArgs(t *testing.T) {
	assert := assert.New(t)
	req := newReq()
	req.Query(&struct {
		Example struct {
			Message string `json:"message"`
			Boolean bool   `json:"boolean"`
			Number  int    `json:"number"`
		} `json:"example" gql_params:"name:Int"`
	}{}).WithVariable("name", 1)
	reqStr := req.GetQuery()
	assert.Equal(`query($name:Int){
    example(name:$name){
        message
        boolean
        number
    }
}
`, reqStr)
}

type mockTransport2 struct {
	calledNum int
}

func (m *mockTransport2) Transport(req Request) (Response, error) {
	m.calledNum++
	return Response{}, nil
}

func TestSendCallsTransport(t *testing.T) {
	assert := assert.New(t)
	r := newReq()
	m := &mockTransport2{}
	r.SetTransport(m)
	r.Send()
	assert.Equal(1, m.calledNum)
}

func TestCanAliasFields(t *testing.T) {
	assert := assert.New(t)
	req := newReq()
	req.Query(&struct {
		Example struct {
			Message string `json:"message" gql:"gql_message"`
		} `json:"example" gql_params:"name:Int"`
	}{}).WithVariable("name", 1)
	reqStr := req.GetQuery()
	assert.Equal(`query($name:Int){
    example(name:$name){
        message: gql_message
    }
}
`, reqStr)
}
func TestCanBuildASimpleQueryWithMultipleArgs(t *testing.T) {
	assert := assert.New(t)
	req := newReq()
	req.Query(&struct {
		Example struct {
			Message string `json:"message"`
			Boolean bool   `json:"boolean"`
			Number  int    `json:"number"`
		} `json:"example" gql_params:"name:Int,id:ID"`
	}{}).WithVariable("name", 1).WithVariable("id", "1234")
	reqStr := req.GetQuery()
	assert.Equal(`query($name:Int, $id:ID){
    example(name:$name, id:$id){
        message
        boolean
        number
    }
}
`, reqStr)
}
