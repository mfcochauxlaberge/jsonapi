package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

var _ Collection = (*WrapperCollection)(nil)

func TestWrapCollection(t *testing.T) {
	assert := assert.New(t)

	res := Wrap(mocktype{})
	col := WrapCollection(res)

	// Collection's type == resource's type
	assert.True(col.GetType().Equal(res.GetType()))

	assert.Equal(0, col.Len())

	// Resource added
	col.Add(res)
	assert.Equal(1, col.Len())
	assert.True(Equal(res, col.At(0)))

	// Index out of bound
	assert.Nil(col.At(999))
}
