package jsonapi

import (
	"testing"

	"kkaribu/tchek"
)

// TestDocument ...
func TestDocument(t *testing.T) {
	pl1 := Document{}

	tchek.AreEqual(t, 0, nil, pl1.Resource)
	tchek.AreEqual(t, 1, nil, pl1.Collection)
	tchek.AreEqual(t, 0, Identifier{}, pl1.Identifier)
	tchek.AreEqual(t, 1, 0, len(pl1.Identifiers))
}
