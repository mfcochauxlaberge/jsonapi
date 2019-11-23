package jsonapi_test

import (
	"encoding/json"
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

func TestUnmarshalIdentifiers(t *testing.T) {
	// Setup
	typ, _ := BuildType(mocktype{})
	typ.NewFunc = func() Resource {
		return Wrap(&mocktype{})
	}
	schema := &Schema{Types: []Type{typ}}

	t.Run("identifier", func(t *testing.T) {
		assert := assert.New(t)

		iden := Identifier{
			ID:   "id2",
			Type: "mocktype",
		}

		payload, err := json.Marshal(iden)
		assert.NoError(err)

		iden2, err := UnmarshalIdentifier(payload, schema)
		assert.NoError(err)
		assert.Equal(iden, iden2)
	})

	t.Run("identifers", func(t *testing.T) {
		assert := assert.New(t)

		idens := Identifiers{
			Identifier{
				ID:   "id2",
				Type: "mocktype",
			},
			Identifier{
				ID:   "id3",
				Type: "mocktype",
			},
		}

		payload, err := json.Marshal(idens)
		assert.NoError(err)

		idens2, err := UnmarshalIdentifiers(payload, schema)
		assert.NoError(err)
		assert.Equal(idens, idens2)
	})
}
