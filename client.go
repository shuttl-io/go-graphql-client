package graphql

import "net/http"

// Request represents a graphql request to some API
type Request interface {

	// SetTransport sets the transport for the request when send
	SetTransport(transport Transport) Request

	//WithVariable adds a variable to the request
	WithVariable(name string, value interface{}) Request

	// Query sets the request to a query. It will take the interface and perform reflection to see
	// what fields to request from graphql. The fields should be json serializable (that this
	// Respects the json tags) or you could use the graphql tags for more specific uses
	Query(object interface{}) Request

	// Mutation sets the request to a mutation. It will take the interface and perform reflection to see
	// what fields to request from graphql. The fields should be json serializable (that this
	// Respects the json tags) or you could use the graphql tags for more specific uses
	Mutation(object interface{}) Request

	// Send sends the request to the graphql API. This returns a filled out version of the interface passed
	// into query or mutation methods
	Send() (Response, error)

	// GetQuery gets the full query
	GetQuery() string

	//GetVariables gets the variables for the request
	GetVariables() map[string]interface{}

	//GetInterface gets the interface that the response should adhere to
	GetInterface() interface{}
}

// Response represents a response from a graphql api
type Response struct {
	// HttpResponse gets the header off the response
	HttpResponse *http.Response
	// Payload is the raw payload
	Payload []byte
	// obj is the object to unserialize
	Response interface{}
	// HttpRequest is the raw Request
	HttpRequest *http.Request
}

// Transport transports the request to an API and returns the response
type Transport interface {
	// ModifyRequest takes a request and then returns the request.
	Transport(req Request) (Response, error)
}

// Client represents a GraphQL client
type Client interface {
	// NewRequest makes a new request to send to some graphql response
	NewRequest() Request

	// AddRequestModifier adds a request modifier to the client, use this to add headers or do some retry logic
	SetTransport(transport Transport) Client
}

type client struct {
	transport Transport
}

// NewClient returns a new Graphql Client
func NewClient(transport Transport) Client {
	return &client{
		transport: transport,
	}
}

func (c *client) NewRequest() Request {
	req := newReq()
	req.SetTransport(c.transport)
	return req
}

func (c *client) SetTransport(t Transport) Client {
	c.transport = t
	return c
}
