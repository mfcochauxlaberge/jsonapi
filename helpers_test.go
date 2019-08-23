package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestBuildType(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		MustBuildType("invalid")
	})

	mock := mockType1{
		ID:    "abc13",
		Str:   "string",
		Int:   -42,
		Uint8: 12,
	}
	typ, err := BuildType(mock)
	assert.NoError(err)
	assert.Equal(true, Equal(Wrap(&mockType1{}), typ.New()))
}
