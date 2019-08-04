package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestReflect(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		MustReflect("invalid")
	})
}
