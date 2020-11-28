package graphql

import "strings"

func buildString(parts []string, level int, builder *strings.Builder, ident int) *strings.Builder {
	if len(parts) == 0 {
		return builder
	}
	part := parts[0]
	newParts := parts[1:]
	if part == "{" {
		builder.WriteString(part)
		return buildString(newParts, level+1, builder, ident)
	}
	if part == "}" {
		level = level - 1
		spaces := strings.Repeat(" ", level*ident)
		builder.WriteString("\n")
		builder.WriteString(spaces)
		builder.WriteString(part)
		return buildString(newParts, level, builder, ident)
	}
	builder.WriteString("\n")
	spaces := strings.Repeat(" ", level*ident)
	builder.WriteString(spaces)
	builder.WriteString(part)
	return buildString(newParts, level, builder, ident)
}

func prettyPrintParts(parts []string, ident int) string {
	return buildString(parts, 0, &strings.Builder{}, ident).String()
}
