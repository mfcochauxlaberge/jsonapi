package jsonapi_test

import (
	"testing"

	"github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestMeta(t *testing.T) {
	assert := assert.New(t)

	meta := jsonapi.Meta{
		"bool":   true,
		"string": "a string",
		"int":    42,
	}

	assert.Equal("a string", meta.GetString("string"))
}
