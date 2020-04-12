package jsonapi_test

import (
	"testing"
	"time"

	"github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestMeta(t *testing.T) {
	assert := assert.New(t)

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
		"time":   "2012-05-16T17:45:28.2539Z",
		"bytes":  []byte{'a', 'b', 'c'},
	}

	assert.True(meta.Has("string"))
	assert.False(meta.Has("unknown"))

	assert.Equal("str", meta.GetString("string"))
	assert.Equal("-12", meta.GetString("int"))
	assert.Equal("-22", meta.GetString("int8"))
	assert.Equal("-32", meta.GetString("int16"))
	assert.Equal("-42", meta.GetString("int32"))
	assert.Equal("-52", meta.GetString("int64"))
	assert.Equal("12", meta.GetString("uint"))
	assert.Equal("22", meta.GetString("uint8"))
	assert.Equal("32", meta.GetString("uint16"))
	assert.Equal("42", meta.GetString("uint32"))
	assert.Equal("52", meta.GetString("uint64"))
	assert.Equal("true", meta.GetString("bool"))
	assert.Equal("2012-05-16T17:45:28.2539Z", meta.GetString("time"))

	assert.Equal(0, meta.GetInt("string"))
	assert.Equal(-12, meta.GetInt("int"))
	assert.Equal(-22, meta.GetInt("int8"))
	assert.Equal(-32, meta.GetInt("int16"))
	assert.Equal(-42, meta.GetInt("int32"))
	assert.Equal(-52, meta.GetInt("int64"))
	assert.Equal(12, meta.GetInt("uint"))
	assert.Equal(22, meta.GetInt("uint8"))
	assert.Equal(32, meta.GetInt("uint16"))
	assert.Equal(42, meta.GetInt("uint32"))
	assert.Equal(52, meta.GetInt("uint64"))
	assert.Equal(0, meta.GetInt("bool"))
	assert.Equal(0, meta.GetInt("time"))

	assert.Equal(false, meta.GetBool("string"))
	assert.Equal(false, meta.GetBool("int"))
	assert.Equal(false, meta.GetBool("int8"))
	assert.Equal(false, meta.GetBool("int16"))
	assert.Equal(false, meta.GetBool("int32"))
	assert.Equal(false, meta.GetBool("int64"))
	assert.Equal(false, meta.GetBool("uint"))
	assert.Equal(false, meta.GetBool("uint8"))
	assert.Equal(false, meta.GetBool("uint16"))
	assert.Equal(false, meta.GetBool("uint32"))
	assert.Equal(false, meta.GetBool("uint64"))
	assert.Equal(true, meta.GetBool("bool"))
	assert.Equal(false, meta.GetBool("time"))

	tm, _ := time.Parse(time.RFC3339Nano, "2012-05-16T17:45:28.2539Z")

	assert.Equal(time.Time{}, meta.GetTime("string"))
	assert.Equal(time.Time{}, meta.GetTime("int"))
	assert.Equal(time.Time{}, meta.GetTime("int8"))
	assert.Equal(time.Time{}, meta.GetTime("int16"))
	assert.Equal(time.Time{}, meta.GetTime("int32"))
	assert.Equal(time.Time{}, meta.GetTime("int64"))
	assert.Equal(time.Time{}, meta.GetTime("uint"))
	assert.Equal(time.Time{}, meta.GetTime("uint8"))
	assert.Equal(time.Time{}, meta.GetTime("uint16"))
	assert.Equal(time.Time{}, meta.GetTime("uint32"))
	assert.Equal(time.Time{}, meta.GetTime("uint64"))
	assert.Equal(time.Time{}, meta.GetTime("bool"))
	assert.Equal(tm, meta.GetTime("time"))
}
