# go-grapqhl-client

A simple graphql library to write and send graphql queries as part of go structs using tags and
easily send and receive graphql results. This takes advantage of the JSON tags to serialize an
object to a graphql query as well as deserialize the response.

## QuickStart

### Installation

A normal go get should do the trick: `go get github.com/shuttl-io/go-graphql-client`.

### Basics

First, you will need to create a new Graphql Client. This is essentially a factory for creating new requests.
The Graphql Client requires a Transport object. The Transport object is what takes the graphql request,
sends it to a graphql api, and then returns the response and deserializes the response. For simplicity, we have
included a simple HTTP transport that allows you to send the API and set headers if you need it:

```golang
package main

import (
    "github.com/shuttl-io/go-graphql-client"
)

func main() {
    client := graphql.NewClient(graphql.NewSimpleHTTPTransport("https://api.example.com/graphql"))
    ...
}
```

Breaking this down, the `graphql.NewClient` takes a `Transport` object and uses that to create and
send requests to some API. The `graphql.NewSimpleHTTPTransport` takes the URL of the api so it knows
where to route requests to.

After that you can use a normal golang struct to make a request against your api:

```golang
package main

import (
    "github.com/shuttl-io/go-graphql-client"
)

type HeroObject struct {
    Name string `json:"name"`
}

type HeroWithFriend struct {
    Name string `json:"name"`
    Friend []HeroObject `json:"friends"`
}

type GraphQLRequest struct {
    Hero Hero `json:"hero"`
}

func main() {
    client := graphql.NewClient(graphql.NewSimpleHTTPTransport("https://api.example.com/graphql"))
    req := &GraphQLRequest{}
    resp, err := client.NewRequest().Query(req).Send()
    fmt.Print(req.Hero.Name)
}
```

This will send a request to your graphql api that looks like this:

```graphql
query {
  hero {
    name
    friends {
      name
    }
  }
}
```

`client.NewRequest()` starts a new Request context that you can use to add a query or mutation and variables.
You then pass your query struct into `.Query()` or `.Mutation()` (depending on what you are trying to do). When
the request comes back from the server, it will be automatically unmarshalled into the `&GraphQLRequest{}`
object on line 21. `.Send()` actually sends the request to the server and will return a Response object that
contains the `http.Response` and `http.Request` objects as well as the raw bytes of the response in `Response.Payload`
as well as the deserialized object on `Request.Response` that is just an `interface{}` which you can then cast to your
graphql request object.

### Adding params

Now we have the basics, how do we query a graphql api? Its really simple in this library. All you need to do is add a tag
to your struct that defines a field as queryable and its params. This tag is the `gql_params` tag and it takes the
following format as a value: `"<argument name>:<graphql type>"`. This tag is comma separated for multiple args After this
is done, then calling `.WithVariable` on the request will format the query correctly:

```golang
package main

import (
    "github.com/shuttl-io/go-graphql-client"
)

type HeroObject struct {
    Name string `json:"name"`
}

type HeroWithFriend struct {
    Name string `json:"name"`
    Friend []HeroObject `json:"friends"`
}

type GraphQLRequest struct {
    Hero Hero `json:"hero" gql_params:"id:ID"`
}

func main() {
    client := graphql.NewClient(graphql.NewSimpleHTTPTransport("https://api.example.com/graphql"))
    req := &GraphQLRequest{}
    resp, err := client.NewRequest()
        .Query(req)
        .WithVariable("id", "1000")
        .Send()
    fmt.Print(req.Hero.Name)
}
```

This will send a query that looks like this:

```graphql
query($id: ID) {
  hero(id: $id) {
    name
    friends {
      name
    }
  }
}
```

With the value for ID being sent in the variables argument of the request. If you don't pass an argument via the
`.WithVariable`, this library will simply not format any arguments on to the request. That is, the request will
only contain variables that the request was asked to include.

### Ignoring a specific field

If you need to ignore a specific field but want it on the query. you can add the tag and value `gql:"omit"` to the
struct. This will not add the field to the query

```golang
package main

import (
    "github.com/shuttl-io/go-graphql-client"
)

type HeroObject struct {
    Name string `json:"name"`
    IgnoreField `json:"ignore" gql:"omit"`
}

type HeroWithFriend struct {
    Name string `json:"name"`
    Friend []HeroObject `json:"friends"`
}

type GraphQLRequest struct {
    Hero Hero `json:"hero" gql_params:"id:ID"`
}

func main() {
    client := graphql.NewClient(graphql.NewSimpleHTTPTransport("https://api.example.com/graphql"))
    req := &GraphQLRequest{}
    resp, err := client.NewRequest()
        .Query(req)
        .WithVariable("id", "1000")
        .Send()
    fmt.Print(req.Hero.Name)
}
```

This will send a query that looks like this:

```graphql
query($id: ID) {
  hero(id: $id) {
    name
    friends {
      name
    }
  }
}
```

### Aliasing field

Sometimes it is necessary to alias a graphql field. To do that, use the JSON tag and the gql tag at the same time:

```golang
package main

import (
    "github.com/shuttl-io/go-graphql-client"
)

type HeroObject struct {
    Name string `json:"hero_name" gql:"name"`
    IgnoreField `json:"ignore" gql:"omit"`
}

type HeroWithFriend struct {
    Name string `json:"name"`
    Friend []HeroObject `json:"friends"`
}

type GraphQLRequest struct {
    Hero Hero `json:"hero" gql_params:"id:ID"`
}

func main() {
    client := graphql.NewClient(graphql.NewSimpleHTTPTransport("https://api.example.com/graphql"))
    req := &GraphQLRequest{}
    resp, err := client.NewRequest()
        .Query(req)
        .WithVariable("id", "1000")
        .Send()
    fmt.Print(req.Hero.Name)
}
```

This will send a query that looks like this:

```graphql
query($id: ID) {
  hero(id: $id) {
    hero_name: name
    friends {
      name
    }
  }
}
```

## Full Working and Copy Pastable Code

```golang
package main

import (
	"fmt"

	"github.com/shuttl-io/go-graphql-client"
)

type continent struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type continentsRequest struct {
	Continents []continent `json:"continents" gql_params:"filter:ContinentFilterInput"`
}

type stringQuery struct {
	Eq    string   `json:"eq,omitempty"`
	Ne    string   `json:"ne,omitempty"`
	In    []string `json:"in,omitempty"`
	Nin   string   `json:"nin,omitempty"`
	Regex string   `json:"regex,omitempty"`
	Glob  string   `json:"glob,omitempty"`
}

type ContinentFilter struct {
	Code stringQuery `json:"code"`
}

func main() {
	client := graphql.NewClient(graphql.NewSimpleHTTPTransport("https://countries.trevorblades.com/"))
	req := &continentsRequest{}
	filter := ContinentFilter{}
	filter.Code.Eq = "AF"
	_, err := client.NewRequest().Query(req).WithVariable("filter", filter).Send()
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	for _, continent := range req.Continents {
		fmt.Println("Continent Code:", continent.Code, "Name:", continent.Name)
		fmt.Println("========================================================")
	}
}
```

## Contributions

Contributions are 100% encouraged. This can take many forms. Just using this project is contribution enough.
If you run into an issue, please drop a line in the issues and we will get back to you ASAP. If you want to
open a PR, that is cool too, just go ahead and open one.
