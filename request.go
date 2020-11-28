package graphql

import "strings"

type request struct {
	tp        string
	retVal    interface{}
	m         *Marshaler
	argValues map[string]interface{}
	err       error
	transport Transport
}

func newReq() *request {
	return &request{
		tp:        "",
		retVal:    nil,
		m:         NewMarshaler(),
		argValues: map[string]interface{}{},
		err:       nil,
	}
}

func (r *request) makeReq(obj interface{}, tp string) Request {
	r.tp = tp
	r.retVal = obj
	_, r.err = r.m.MarshalToGraphql(obj)
	return r
}

func (r *request) Query(object interface{}) Request {
	return r.makeReq(object, "query")
}

func (r *request) Mutation(object interface{}) Request {
	return r.makeReq(object, "mutation")
}

func (r *request) WithVariable(name string, value interface{}) Request {
	r.argValues[name] = value
	return r
}

func (r *request) SetTransport(transport Transport) Request {
	r.transport = transport
	return r
}

func (r *request) GetInterface() interface{} {
	return r.retVal
}

func (r *request) Send() (Response, error) {
	return r.transport.Transport(r)
}

func (r *request) GetVariables() map[string]interface{} {
	return r.argValues
}

func (r *request) GetQuery() string {
	for key := range r.argValues {
		r.m.rootPart.markArgAsNeeded(key)
	}
	collectedArgs := r.m.rootPart.collectArgs()
	builders := &strings.Builder{}
	builders.WriteString(r.tp)
	if len(collectedArgs) > 0 {
		builders.WriteString("(")
		builders.WriteString(strings.Join(collectedArgs, ", "))
		builders.WriteString(")")
	}
	r.m.rootPart.buildStr(builders, 0)
	return builders.String()
}
