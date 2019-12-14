package jsonapi_test

import (
	"testing"
	"time"

	"github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestMeta(t *testing.T) {
	assert := assert.New(t)

	tm, _ := time.Parse(time.RFC3339Nano, "2012-05-16T17:45:28.2539Z")

	meta := jsonapi.Meta{
		"string": "str",
		"int":    -12,
		"int8":   -22,
		"int16":  -32,
		"int32":  -42,
		"int64":  -52,
		"uint":   12,
		"uint8":  22,
		"uint16": 32,
		"uint32": 42,
		"uint64": 52,
		"bool":   true,
		"time":   tm,
		"bytes":  []byte{'a', 'b', 'c'},
	}

	assert.Equal("str", meta.GetString("string"))
	assert.Equal("-12", meta.GetString("int"))
	assert.Equal("-22", meta.GetString("int8"))
	assert.Equal("-32", meta.GetString("int16"))
	assert.Equal("-42", meta.GetString("int32"))
	assert.Equal("-52", meta.GetString("int64"))
	assert.Equal(" 12", meta.GetString("uint"))
	assert.Equal(" 22", meta.GetString("uint8"))
	assert.Equal(" 32", meta.GetString("uint16"))
	assert.Equal(" 42", meta.GetString("uint32"))
	assert.Equal(" 52", meta.GetString("uint64"))
}
