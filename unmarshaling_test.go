package jsonapi

import (
	"testing"

	"github.com/mfcochauxlaberge/tchek"
)

func TestUnmarshalResource(t *testing.T) {
	schema := NewMockSchema()

	res1 := Wrap(&MockType3{
		ID:    "mt1",
		Attr1: "a string",
		Attr2: 1,
		Rel1:  "mt2",
		Rel2:  []string{"mt3", "mt4"},
	})

	url1, err := ParseRawURL(schema, "/mocktypes3/mt1")
	tchek.UnintendedError(err)

	meta1 := map[string]interface{}{
		"str": "a string\\^ç\"",
		"num": float64(42),
		"b":   true,
	}

	doc1 := NewDocument()
	doc1.Data = res1
	doc1.Meta = meta1

	body1, err := Marshal(doc1, url1)
	tchek.UnintendedError(err)

	pl1, err := Unmarshal(body1, url1, schema)
	tchek.UnintendedError(err)

	dst1 := pl1.Data.(Resource)

	tchek.HaveEqualAttributes(t, "same attribues", res1, dst1)
	tchek.AreEqual(t, "same meta object", meta1, pl1.Meta)
}

func TestUnmarshalIdentifier(t *testing.T) {
	schema := NewMockSchema()

	id1 := Identifier{ID: "abc123", Type: "mocktypes1"}

	url1, err := ParseRawURL(schema, "/mocktypes3/mt1/relationships/rel1")
	tchek.UnintendedError(err)

	meta1 := map[string]interface{}{
		"str": "a string\\^ç\"",
		"num": float64(42),
		"b":   true,
	}

	doc1 := NewDocument()
	doc1.Data = id1
	doc1.Meta = meta1

	body1, err := Marshal(doc1, url1)
	tchek.UnintendedError(err)

	pl1, err := Unmarshal(body1, url1, schema)
	tchek.UnintendedError(err)

	dst1 := pl1.Data.(Identifier)

	tchek.AreEqual(t, "same identifier", id1, dst1)
	tchek.AreEqual(t, "same meta map", meta1, pl1.Meta)
}

func TestUnmarshalIdentifiers(t *testing.T) {
	schema := NewMockSchema()

	ids1 := Identifiers{
		Identifier{ID: "abc123", Type: "mocktypes1"},
		Identifier{ID: "def456", Type: "mocktypes1"},
		Identifier{ID: "ghi789", Type: "mocktypes1"},
	}

	url1, err := ParseRawURL(schema, "/mocktypes3/mt1/relationships/rel2")
	tchek.UnintendedError(err)

	meta1 := map[string]interface{}{
		"str": "a string\\^ç\"",
		"num": float64(42),
		"b":   true,
	}

	doc1 := NewDocument()
	doc1.Data = ids1
	doc1.Meta = meta1

	body1, err := Marshal(doc1, url1)
	tchek.UnintendedError(err)

	pl1, err := Unmarshal(body1, url1, schema)
	tchek.UnintendedError(err)

	dst1 := pl1.Data.(Identifiers)

	tchek.AreEqual(t, "same identifiers", ids1, dst1)
	tchek.AreEqual(t, "same meta map", meta1, pl1.Meta)
}
