package graphql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func noSpaces(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

func TestCanConvertBasicToGql(t *testing.T) {
	realQ := `{
  example{
	message
	boolean
	number
  }
}
`
	marshaler := NewMarshaler()
	marshaler.IdentLevel = 2
	assert := assert.New(t)
	q, err := marshaler.MarshalToGraphql(&struct {
		Example struct {
			Message string `json:"message"`
			Boolean bool   `json:"boolean"`
			Number  int    `json:"number"`
		} `json:"example" gql_params:"name:Int"`
	}{})
	assert.NoError(err)
	assert.Equal(noSpaces(realQ), noSpaces(q))
}

type someMock struct {
}

func (s someMock) MarshalGql(m *Marshaler) ([]*QueryPart, error) {
	qp := NewQueryPart("someMockField1")
	qp.SubFields = append(qp.SubFields,
		NewQueryPart("someDeepMock"),
		NewQueryPart("someotherDeepMock"),
	)
	return []*QueryPart{
		qp,
		NewQueryPart("someMockField2"),
	}, nil
}

func TestCanConvertMarshaler(t *testing.T) {
	realQ := `{
  example{
    someMock{
      someMockField1{
        someDeepMock
        someotherDeepMock
      }
    someMockField2
  }
}
}
`
	marshaler := NewMarshaler()
	marshaler.IdentLevel = 2
	assert := assert.New(t)
	q, err := marshaler.MarshalToGraphql(&struct {
		Example struct {
			SomeMock someMock `json:"someMock"`
		} `json:"example" gql_params:"name:Int"`
	}{})
	assert.NoError(err)
	assert.Equal(noSpaces(realQ), noSpaces(q))
}

func TestCanOmit(t *testing.T) {
	realQ := `{
  example{
	message
	boolean
	number
    SomeValue
  }
}
`
	marshaler := NewMarshaler()
	marshaler.IdentLevel = 2
	assert := assert.New(t)
	q, err := marshaler.MarshalToGraphql(&struct {
		Example struct {
			Message   string `json:"message"`
			Boolean   bool   `json:"boolean"`
			Number    int    `json:"number"`
			Omit      string `json:"will_omit_in_gql" gql:"omit"`
			SomeValue string
		} `json:"example" gql_params:"name:Int"`
	}{})
	assert.NoError(err)
	assert.Equal(noSpaces(realQ), noSpaces(q))
}

func TestCanConvertToString(t *testing.T) {
	assert := assert.New(t)
	marshaler := NewMarshaler()
	marshaler.MarshalToGraphql(&struct {
		Example struct {
			Message string `gql:"message" json:"thisValueIsForJsonOnly"`
			Boolean bool   `json:"boolean"`
			Number  int    `json:"number"`
			Omit    string `json:"will_omit_in_gql" gql:"omit"`
		} `json:"example" gql_params:"name:Int"`
	}{})
	assert.NotEmpty(marshaler.String())
}

func TestCanMarshal(t *testing.T) {
	assert := assert.New(t)
	m, err := Marshal(&struct {
		Example struct {
			Message string `gql:"message" json:"thisValueIsForJsonOnly"`
			Boolean bool   `json:"boolean"`
			Number  int    `json:"number"`
			Omit    string `json:"will_omit_in_gql" gql:"omit"`
		} `json:"example" gql_params:"name:Int"`
	}{})
	assert.NoError(err)
	assert.NotNil(m)
}

func TestMakeArgRequired(t *testing.T) {
	assert := assert.New(t)
	m, _ := Marshal(&struct {
		Example struct {
			Message string `gql:"message" json:"thisValueIsForJsonOnly"`
			Boolean bool   `json:"boolean"`
			Number  int    `json:"number"`
			Omit    string `json:"will_omit_in_gql" gql:"omit"`
		} `json:"example" gql_params:"name:Int, foo:String"`
	}{})
	m.AddToArgs("name", "foo")
	assert.Len(m.rootPart.collectArgs(), 2)
}

func TestMakeArgRequiredFromMarshal(t *testing.T) {
	assert := assert.New(t)
	marshaler := NewMarshaler()
	marshaler.MarshalToGraphql(&struct {
		Example struct {
			Message string `gql:"message" json:"thisValueIsForJsonOnly"`
			Boolean bool   `json:"boolean"`
			Number  int    `json:"number"`
			Omit    string `json:"will_omit_in_gql" gql:"omit"`
		} `json:"example" gql_params:"name:Int, foo:String"`
	}{}, "name", "foo")
	assert.Len(marshaler.rootPart.collectArgs(), 2)
}

type testExample struct {
	Message string `json:"tst"`
}

func TestCanConvertArray(t *testing.T) {
	realQ := `{
  test{
	example{
		tst
	}
  }
}
`
	marshaler := NewMarshaler()
	marshaler.IdentLevel = 2
	assert := assert.New(t)
	q, err := marshaler.MarshalToGraphql(&struct {
		Example struct {
			Test []testExample `json:"example" gql_params:"name:Int"`
		} `json:"test"`
	}{})
	assert.NoError(err)
	assert.Equal(noSpaces(realQ), noSpaces(q))
}

func TestCanConvertMarshalerSlice(t *testing.T) {
	realQ := `{
  example{
    someMock{
      someMockField1{
        someDeepMock
        someotherDeepMock
      }
    someMockField2
  }
}
}
`
	marshaler := NewMarshaler()
	marshaler.IdentLevel = 2
	assert := assert.New(t)
	q, err := marshaler.MarshalToGraphql(&struct {
		Example struct {
			SomeMock []someMock `json:"someMock"`
		} `json:"example" gql_params:"name:Int"`
	}{})
	assert.NoError(err)
	assert.Equal(noSpaces(realQ), noSpaces(q))
}
