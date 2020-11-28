package graphql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Marshaler marshels the object to a graphql request
type Marshaler struct {
	argumentTypes map[string]string
	rootPart      *QueryPart
	includedArgs  map[string]bool
	IdentLevel    int
}

// GqlMarshaler implements a way to Marshal to a graphql request. This returns a list of QueryParts
// to add to the parent part of the graphql API.
type GqlMarshaler interface {
	MarshalGql(marshaler *Marshaler) ([]*QueryPart, error)
}

// NewMarshaler returns a new Marshaler
func NewMarshaler() *Marshaler {
	return &Marshaler{
		argumentTypes: make(map[string]string),
		rootPart:      nil,
		includedArgs:  map[string]bool{},
		IdentLevel:    0,
	}
}

func (g *Marshaler) marshalStruct(tp reflect.Type, val reflect.Value, rootPart *QueryPart) error {
	numFields := tp.NumField()
	for i := 0; i < numFields; i++ {
		field := tp.Field(i)
		gqlTags, hasGqlTags := field.Tag.Lookup("gql")
		tags := newSet()
		if hasGqlTags {
			tags = newSet(strings.Split(gqlTags, ",")...)
		}
		gqlParams, hasParams := field.Tag.Lookup("gql_params")
		params := newSet(strings.Split(gqlParams, ",")...)
		json, hasJson := field.Tag.Lookup("json")
		values := strings.Split(gqlTags, ",")
		jsonValues := strings.Split(json, ",")
		if tags.has("omit") {
			continue
		}
		var name string
		if !hasJson {
			name = field.Name
		} else {
			name = jsonValues[0]
		}
		if hasJson && hasGqlTags {
			name = fmt.Sprintf("%s: %s", jsonValues[0], values[0])
		}
		part := NewQueryPart(name)
		if hasParams {
			for _, param := range params.elements() {
				parts := strings.Split(strings.Trim(param, " "), ":")
				name := parts[0]
				tps := parts[1]
				part.Arguments[name] = tps
			}
		}

		g.marshal(tp.Field(i).Type, val.Field(i), part)
		rootPart.SubFields = append(rootPart.SubFields, part)
	}
	return nil
}

func (g *Marshaler) marshal(tp reflect.Type, val reflect.Value, parent *QueryPart) error {
	mType := reflect.TypeOf(new(GqlMarshaler)).Elem()
	if tp.Implements(mType) {
		m := val.Interface().(GqlMarshaler)
		queryParties, err := m.MarshalGql(g)
		if err != nil {
			return nil
		}
		parent.SubFields = append(parent.SubFields, queryParties...)
		return nil
	}
	switch kind := tp.Kind(); kind {
	case reflect.Ptr:
		if val.IsNil() {
			val = reflect.New(tp.Elem())
		}
		return g.marshal(tp.Elem(), val.Elem(), parent)
	case reflect.Struct:
		return g.marshalStruct(tp, val, parent)
	case reflect.String:
		return nil
	case reflect.Slice:
		return g.marshal(tp.Elem(), reflect.New(tp.Elem()).Elem(), parent)
	default:
		return fmt.Errorf("unrecognized kind: %s. this isn't supposed to happen, let us know if it does", kind)
	}
}

// MarshalToGraphql takes the object and returns a new graphql query
func (g *Marshaler) MarshalToGraphql(obj interface{}, argsToInclude ...string) (string, error) {
	tp := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	g.rootPart = NewQueryPart("")
	if tp.Kind() != reflect.Ptr {
		return "", errors.New("object should be a pointer")
	}
	err := g.marshal(tp, val, g.rootPart)
	for _, arg := range argsToInclude {
		g.rootPart.markArgAsNeeded(arg)
	}
	return g.joinParts(), err
}

func (g *Marshaler) joinParts() string {
	return g.rootPart.String()
}

// String converts the current state of the Marshaler to a string
func (g *Marshaler) String() string {
	return g.rootPart.String()
}

// AddToArgs adds a new arg name to be added as part of the query
func (g *Marshaler) AddToArgs(argNames ...string) {
	for _, argName := range argNames {
		g.rootPart.markArgAsNeeded(argName)
	}
}

// Marshal will start the marshalling process and convert the object to a graphql request
func Marshal(obj interface{}) (*Marshaler, error) {
	m := NewMarshaler()
	_, err := m.MarshalToGraphql(obj)
	return m, err
}
