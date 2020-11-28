package graphql

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testQuery struct {
	Message string `json:"message"`
	SubQ    struct {
		SubMessage string `json:"message"`
	} `json:"sub_query" gql_params:"name:String"`
}

type testRep struct {
	Data struct {
		Message string `json:"message"`
		Sub     struct {
			Message string `json:"message"`
		} `json:"sub_query"`
	} `json:"data"`
	Error []gqlError `json:"errors"`
}

type reqObj struct {
	Query     string            `json:"query"`
	Variables map[string]string `json:"variables"`
}

func TestCallsTheRightUrlCorrectly(t *testing.T) {
	assert := assert.New(t)
	wasCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		wasCalled = true
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	transport := NewSimpleHTTPTransport(server.URL)
	assert.NotNil(transport)
	msg := &testQuery{}
	req := newReq().Query(msg)
	transport.Transport(req)
	assert.True(wasCalled)
}

func TestSendsTheRightRequest(t *testing.T) {
	assert := assert.New(t)
	wasCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		buf := bytes.Buffer{}
		defer req.Body.Close()
		buf.ReadFrom(req.Body)
		d := reqObj{}
		err := json.Unmarshal(buf.Bytes(), &d)
		assert.NoError(err)
		assert.NotEmpty(d.Query)
		assert.Equal(noSpaces(`query{
			message
			sub_query {
				message
			}
		}
		`), noSpaces(d.Query))
		// Test request parameters
		wasCalled = true

		// Send response to be tested
		resp := testRep{}
		resp.Data.Message = "Good"
		resp.Data.Sub.Message = "Great"
		data, _ := json.Marshal(resp)
		rw.Write(data)
	}))
	defer server.Close()
	transport := NewSimpleHTTPTransport(server.URL)
	assert.NotNil(transport)
	msg := &testQuery{}
	req := newReq().Query(msg)
	resp, err := transport.Transport(req)
	assert.NoError(err)
	assert.True(wasCalled)
	assert.Equal("Good", msg.Message)
	assert.Equal("Great", msg.SubQ.SubMessage)
	assert.NotNil(resp.HttpRequest)
	assert.NotNil(resp.HttpResponse)
	assert.NotNil(resp.Payload)
	casted, ok := resp.Response.(*testQuery)
	assert.True(ok)
	assert.Equal("Good", casted.Message)
	assert.Equal("Great", casted.SubQ.SubMessage)
}

func TestCanAddHeader(t *testing.T) {
	assert := assert.New(t)
	wasCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal("application/json", req.Header.Get("Content-Type"))
		wasCalled = true
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	defer server.Close()
	transport := NewSimpleHTTPTransport(server.URL)
	transport.AddHeader("Content-Type", "application/json")
	assert.NotNil(transport)
	msg := &testQuery{}
	req := newReq().Query(msg)
	transport.Transport(req)
	assert.True(wasCalled)
}

func TestSendsTheRightRequestWithVariables(t *testing.T) {
	assert := assert.New(t)
	wasCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		buf := bytes.Buffer{}
		defer req.Body.Close()
		buf.ReadFrom(req.Body)
		d := reqObj{}
		err := json.Unmarshal(buf.Bytes(), &d)
		assert.NoError(err)
		assert.NotEmpty(d.Query)
		assert.Equal(noSpaces(`query($name:String){
			message
			sub_query(name:$name) {
				message
			}
		}
		`), noSpaces(d.Query))

		value, hasValue := d.Variables["name"]
		assert.True(hasValue)
		assert.Equal("someName", value)
		// Test request parameters
		wasCalled = true

		// Send response to be tested
		resp := testRep{}
		resp.Data.Message = "Good"
		resp.Data.Sub.Message = "Great"
		data, _ := json.Marshal(resp)
		rw.Write(data)
	}))
	defer server.Close()
	transport := NewSimpleHTTPTransport(server.URL)
	assert.NotNil(transport)
	msg := &testQuery{}
	req := newReq().Query(msg).WithVariable("name", "someName")
	resp, err := transport.Transport(req)
	assert.NoError(err)
	assert.True(wasCalled)
	assert.Equal("Good", msg.Message)
	assert.Equal("Great", msg.SubQ.SubMessage)
	assert.NotNil(resp.HttpRequest)
	assert.NotNil(resp.HttpResponse)
	assert.NotNil(resp.Payload)
	casted, ok := resp.Response.(*testQuery)
	assert.True(ok)
	assert.Equal("Good", casted.Message)
	assert.Equal("Great", casted.SubQ.SubMessage)
}

func TestHandlesErrorsRight(t *testing.T) {
	assert := assert.New(t)
	wasCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		wasCalled = true

		// Send response to be tested
		resp := testRep{}
		resp.Error = []gqlError{
			{
				Message: "some_error",
			},
		}
		data, _ := json.Marshal(resp)
		rw.Write(data)
	}))
	defer server.Close()
	transport := NewSimpleHTTPTransport(server.URL)
	assert.NotNil(transport)
	msg := &testQuery{}
	req := newReq().Query(msg).WithVariable("name", "someName")
	_, err := transport.Transport(req)
	assert.Error(err)
	assert.Equal("graphql: some_error", err.Error())
	assert.True(wasCalled)
}
