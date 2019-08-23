package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestIdentifiers(t *testing.T) {
	assert := assert.New(t)

	idens := NewIdentifiers("type1", nil)
	assert.Empty(idens)
	assert.Empty(idens.IDs())

	idens = NewIdentifiers("type1", []string{"id1", "id2", "id3"})
	assert.Len(idens, 3)
	assert.Equal(Identifier{ID: "id1", Type: "type1"}, idens[0])
	assert.Equal(Identifier{ID: "id2", Type: "type1"}, idens[1])
	assert.Equal(Identifier{ID: "id3", Type: "type1"}, idens[2])
	assert.Equal([]string{"id1", "id2", "id3"}, idens.IDs())
}
