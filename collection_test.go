package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

var _ Collection = (*Resources)(nil)

func TestResources(t *testing.T) {
	assert := assert.New(t)

	col := Resources{}
	assert.Equal(Type{}, col.GetType())
	assert.Equal(0, col.Len())

	// Add a resource
	res := &SoftResource{}
	res.SetID("id")
	col.Add(res)
	assert.Equal(1, col.Len())

	// Retrieve a resource
	assert.Equal("id", col.At(0).GetID())
	assert.Nil(col.At(1))
}

func TestUnmarshalCollection(t *testing.T) {
	assert := assert.New(t)

	// Invalid payload
	payload := `{"no:valid"}`

	col, err := UnmarshalCollection([]byte(payload), nil)

	assert.Error(err)
	assert.Nil(col)
}
