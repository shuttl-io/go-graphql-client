package graphql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanFormatCorrectly(t *testing.T) {
	assert := assert.New(t)
	parts := []string{
		"{",
		"something",
		"somethingElse", "{",
		"example",
		"}",
		"}",
	}

	builtString := prettyPrintParts(parts, 2)
	fmt.Println(builtString)
	assert.Equal(`{
  something
  somethingElse{
    example
  }
}`, builtString)
}
