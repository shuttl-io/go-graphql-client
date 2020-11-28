package graphql

import (
	"fmt"
	"strings"
)

// QueryPart Represents a part of a graphql query
type QueryPart struct {
	// Value is the name of the field that graphql needs
	Value string
	// Arguments are the arguments that this field can take this maps args to the Graphql type
	Arguments map[string]string
	// Subfield represents other fields, for example in a bigger queries
	SubFields    []*QueryPart
	requiredArgs []string
}

// NewQueryPart creates a new QueryPart
func NewQueryPart(name string) *QueryPart {
	return &QueryPart{
		Value:        name,
		Arguments:    make(map[string]string),
		SubFields:    []*QueryPart{},
		requiredArgs: []string{},
	}
}

func (q *QueryPart) markArgAsNeeded(arg string) {
	if _, ok := q.Arguments[arg]; ok {
		q.requiredArgs = append(q.requiredArgs, arg)
	} else if len(q.SubFields) > 0 {
		for _, sub := range q.SubFields {
			sub.markArgAsNeeded(arg)
		}
	}

}

func (q *QueryPart) collectArgs() []string {
	val := []string{}
	for _, arg := range q.requiredArgs {
		val = append(val, fmt.Sprintf("$%s:%s", arg, q.Arguments[arg]))
	}
	if len(q.SubFields) > 0 {
		for _, sub := range q.SubFields {
			val = append(val, sub.collectArgs()...)
		}
	}
	return val
}

func (q *QueryPart) buildStr(builder *strings.Builder, level int) *strings.Builder {
	builder.WriteString(strings.Repeat(" ", level*4))
	builder.WriteString(q.Value)
	if len(q.requiredArgs) > 0 {
		builder.WriteString("(")
		str := []string{}
		for _, arg := range q.requiredArgs {
			str = append(str, fmt.Sprintf("%s:$%s", arg, arg))
		}
		builder.WriteString(strings.Join(str, ", "))
		builder.WriteString(")")
	}
	if len(q.SubFields) > 0 {
		builder.WriteString("{\n")
		for _, sub := range q.SubFields {
			sub.buildStr(builder, level+1)
		}
		builder.WriteString(strings.Repeat(" ", level*4))
		builder.WriteString("}")
	}
	builder.WriteString("\n")
	return builder
}

func (q *QueryPart) String() string {
	builder := &strings.Builder{}
	q.buildStr(builder, 0)
	return builder.String()
}
