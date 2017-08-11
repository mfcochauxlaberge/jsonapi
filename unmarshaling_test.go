package jsonapi

import (
	"testing"

	"github.com/kkaribu/tchek"
)

func TestUnmarshalResource(t *testing.T) {
	reg := NewMockRegistry()

	res1 := Wrap(&MockType3{
		ID:    "mt1",
		Attr1: "a string",
		Attr2: 1,
		Rel1:  "mt2",
		Rel2:  []string{"mt3", "mt4"},
	})

	url1, err := ParseRawURL(reg, "/mocktypes3/mt1")
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

	pl1, err := Unmarshal(body1, url1, reg)
	tchek.UnintendedError(err)

	dst1 := pl1.Data.(Resource)

	tchek.HaveEqualAttributes(t, -1, res1, dst1)
	tchek.AreEqual(t, -1, meta1, pl1.Meta)
}

func TestUnmarshalCollection(t *testing.T) {
	reg := NewMockRegistry()

	col1 := WrapCollection(Wrap(&MockType1{}))
	url1, err := ParseRawURL(reg, "/mocktypes1")
	tchek.UnintendedError(err)

	meta1 := map[string]interface{}{
		"str": "a string\\^ç\"",
		"num": float64(42),
		"b":   true,
	}

	doc1 := NewDocument()
	doc1.Data = col1
	doc1.Meta = meta1

	body1, err := Marshal(doc1, url1)
	tchek.UnintendedError(err)

	dst1 := WrapCollection(Wrap(&MockType1{}))

	pl1, err := Unmarshal(body1, url1, reg)
	tchek.UnintendedError(err)

	tchek.HaveEqualAttributes(t, -1, col1, dst1)
	tchek.AreEqual(t, -1, meta1, pl1.Meta)
}
