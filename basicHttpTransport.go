package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type gqlError struct {
	Message string `json:"message"`
}

func (g gqlError) Error() string {
	return "graphql: " + g.Message
}

type graphqlResponse struct {
	Data   interface{} `json:"data"`
	Errors []gqlError  `json:"errors"`
}

// SimpleHTTPTransport is a simple http api that allows you to make a single request to headers
// get the response from the API.
type SimpleHTTPTransport struct {
	apiURL  string
	headers map[string][]string
}

// NewSimpleHTTPTransport takes the api URL and then returns a SimpleHttpTransport
func NewSimpleHTTPTransport(apiURL string) *SimpleHTTPTransport {
	transport := &SimpleHTTPTransport{
		apiURL:  apiURL,
		headers: map[string][]string{},
	}
	transport.AddHeader("Content-Type", "application/json")
	return transport
}

// AddHeader adds a header onto therequest object
func (s *SimpleHTTPTransport) AddHeader(name string, value string) {
	headerVals, hasHeader := s.headers[name]
	if !hasHeader {
		headerVals = []string{}
	}
	headerVals = append(headerVals, value)
	s.headers[name] = headerVals
}

// Transport the request to the API
func (s *SimpleHTTPTransport) Transport(req Request) (Response, error) {
	// for  key, value := range req.GetVariables() {
	// 	gqlReq.Var(key, value)
	// }
	response := Response{}
	bts, err := json.Marshal(struct {
		Query    string                 `json:"query"`
		Variable map[string]interface{} `json:"variables"`
	}{
		Query:    req.GetQuery(),
		Variable: req.GetVariables(),
	})
	if err != nil {
		return response, err
	}
	httpReq, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(bts))
	for key, value := range s.headers {
		for _, headerVal := range value {
			httpReq.Header.Add(key, headerVal)
		}
	}
	response.HttpRequest = httpReq
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return response, fmt.Errorf("error from the api: %s", resp.Status)
	}
	response.HttpResponse = resp
	bf := bytes.Buffer{}
	if _, err = bf.ReadFrom(resp.Body); err != nil {
		return response, err
	}
	response.Payload = bf.Bytes()
	data := &graphqlResponse{
		Data: req.GetInterface(),
	}
	err = json.Unmarshal(response.Payload, data)
	if err != nil {
		return response, err
	}
	if len(data.Errors) > 0 {
		return response, data.Errors[0]
	}
	response.Response = data.Data
	return response, nil
}
