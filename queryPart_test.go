package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanConvertBasicQueryPartToString(t *testing.T) {
	assert := assert.New(t)
	qp := NewQueryPart("some_string")
	assert.Equal("some_string\n", qp.String())
}

func TestCanMakeSubfieldsRight(t *testing.T) {
	assert := assert.New(t)
	qp := NewQueryPart("some_string")
	qp.SubFields = append(qp.SubFields,
		NewQueryPart("some_field_one"),
		NewQueryPart("some_field_two"),
		NewQueryPart("some_field_three"),
	)
	assert.Equal(`some_string{
    some_field_one
    some_field_two
    some_field_three
}
`, qp.String())
	qp2 := NewQueryPart("some_field_three")
	qp2.SubFields = append(qp2.SubFields,
		NewQueryPart("sub_sub_one"),
		NewQueryPart("sub_sub_two"),
		NewQueryPart("sub_sub_three"),
	)
	qp = NewQueryPart("some_string")
	qp.SubFields = append(qp.SubFields,
		NewQueryPart("some_field_one"),
		NewQueryPart("some_field_two"),
		qp2,
	)
	assert.Equal(`some_string{
    some_field_one
    some_field_two
    some_field_three{
        sub_sub_one
        sub_sub_two
        sub_sub_three
    }
}
`, qp.String())
}

func TestCanStringifyWithArgs(t *testing.T) {
	assert := assert.New(t)
	qp := NewQueryPart("some_string")
	qp.markArgAsNeeded("some_arg")
	qp.markArgAsNeeded("some_arg2")
	qp.requiredArgs = []string{"some_arg", "some_arg2"}
	assert.Equal("some_string(some_arg:$some_arg, some_arg2:$some_arg2)\n", qp.String())
}

func TestCanStringifySubfieldsWithArgs(t *testing.T) {
	assert := assert.New(t)
	qp := NewQueryPart("some_string")
	qp.Arguments["some_arg"] = ""
	qp.Arguments["some_arg2"] = ""
	qp.markArgAsNeeded("some_arg")
	qp.markArgAsNeeded("some_arg2")
	qp.SubFields = append(qp.SubFields,
		NewQueryPart("some_field_one"),
		NewQueryPart("some_field_two"),
		NewQueryPart("some_field_three"),
	)
	assert.Equal(`some_string(some_arg:$some_arg, some_arg2:$some_arg2){
    some_field_one
    some_field_two
    some_field_three
}
`, qp.String())
}

func TestCanCollectRequiredArgs(t *testing.T) {
	assert := assert.New(t)
	qp2 := NewQueryPart("some_field_three")
	qp2.SubFields = append(qp2.SubFields,
		NewQueryPart("sub_sub_one"),
		NewQueryPart("sub_sub_two"),
		NewQueryPart("sub_sub_three"),
	)
	qp := NewQueryPart("some_string")
	qp.Arguments["some_args"] = "String"
	qp.Arguments["some_arg2"] = "String"
	qp2.Arguments["some_int"] = "Int!"
	qp2.Arguments["some_bool"] = "Bool"
	qp.SubFields = append(qp.SubFields,
		NewQueryPart("some_field_one"),
		NewQueryPart("some_field_two"),
		qp2,
	)
	qp.markArgAsNeeded("some_args")
	qp.markArgAsNeeded("some_int")
	qp.markArgAsNeeded("some_bool")

	args := qp.collectArgs()
	assert.Len(args, 3)
	assert.Equal([]string{"$some_args:String", "$some_int:Int!", "$some_bool:Bool"}, args)
}

func TestAllTogetherNow(t *testing.T) {
	assert := assert.New(t)
	qp2 := NewQueryPart("some_field_three")
	qp2.SubFields = append(qp2.SubFields,
		NewQueryPart("sub_sub_one"),
		NewQueryPart("sub_sub_two"),
		NewQueryPart("sub_sub_three"),
	)
	qp := NewQueryPart("some_string")
	qp.Arguments["some_args"] = "String"
	qp.Arguments["some_arg2"] = "String"
	qp2.Arguments["some_int"] = "Int!"
	qp2.Arguments["some_bool"] = "Bool"
	qp.SubFields = append(qp.SubFields,
		NewQueryPart("some_field_one"),
		NewQueryPart("some_field_two"),
		qp2,
	)
	qp.markArgAsNeeded("some_args")
	qp.markArgAsNeeded("some_int")
	qp.markArgAsNeeded("some_bool")

	assert.Equal(`some_string(some_args:$some_args){
    some_field_one
    some_field_two
    some_field_three(some_int:$some_int, some_bool:$some_bool){
        sub_sub_one
        sub_sub_two
        sub_sub_three
    }
}
`, qp.String())

}
