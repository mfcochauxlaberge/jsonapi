package jsonapi_test

import (
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalResource(t *testing.T) {
	assert := assert.New(t)

	schema := newMockSchema()

	res1 := Wrap(&mockType3{
		ID:    "mt1",
		Attr1: "a string",
		Attr2: 1,
		Rel1:  "mt2",
		Rel2:  []string{"mt3", "mt4"},
	})

	url1, err := NewURLFromRaw(schema, "/mocktypes3/mt1")
	assert.NoError(err)
	meta1 := map[string]interface{}{
		"str": "a string\\^ç\"",
		"num": float64(42),
		"b":   true,
	}

	doc1 := NewDocument()
	doc1.Data = res1
	doc1.Meta = meta1

	body1, err := Marshal(doc1, url1)
	assert.NoError(err)
	pl1, err := Unmarshal(body1, url1, schema)
	assert.NoError(err)
	// dst1 := pl1.Data.(Resource)

	// assert.HaveEqualAttributes(t, "same attribues", res1, dst1) TODO Fix test
	assert.Equal(meta1, pl1.Meta, "same meta object")
}

func TestUnmarshalIdentifier(t *testing.T) {
	assert := assert.New(t)

	schema := newMockSchema()

	id1 := Identifier{ID: "abc123", Type: "mocktypes1"}

	url1, err := NewURLFromRaw(schema, "/mocktypes3/mt1/relationships/rel1")
	assert.NoError(err)
	meta1 := map[string]interface{}{
		"str": "a string\\^ç\"",
		"num": float64(42),
		"b":   true,
	}

	doc1 := NewDocument()
	doc1.Data = id1
	doc1.Meta = meta1

	body1, err := Marshal(doc1, url1)
	assert.NoError(err)
	pl1, err := Unmarshal(body1, url1, schema)
	assert.NoError(err)
	dst1 := pl1.Data.(Identifier)

	assert.Equal(id1, dst1, "same identifier")
	assert.Equal(meta1, pl1.Meta, "same meta map")
}

// func TestUnmarshalIdentifiers(t *testing.T) {
// 	assert := assert.New(t)

// 	schema := newMockSchema()

// 	ids1 := Identifiers{
// 		Identifier{ID: "abc123", Type: "mocktypes1"},
// 		Identifier{ID: "def456", Type: "mocktypes1"},
// 		Identifier{ID: "ghi789", Type: "mocktypes1"},
// 	}

// 	url1, err := NewURLFromRaw(schema, "/mocktypes3/mt1/relationships/rel2")
// 	assert.NoError(err)

// 	meta1 := map[string]interface{}{
// 		"str": "a string\\^ç\"",
// 		"num": float64(42),
// 		"b":   true,
// 	}

// 	doc1 := NewDocument()
// 	doc1.Data = ids1
// 	doc1.Meta = meta1

// 	body1, err := Marshal(doc1, url1)
// 	assert.NoError(err)

// 	pl1, err := Unmarshal(body1, url1, schema)
// 	assert.NoError(err)

// 	dst1 := pl1.Data.(Identifiers)

// 	assert.Equal(ids1, dst1, "same identifiers")
// 	assert.Equal(meta1, pl1.Meta, "same meta map")
// }
