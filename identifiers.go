package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
)

// NewIdentifiers returns an Identifiers object.
//
// t is the type of the identifiers. ids is the set of IDs.
func NewIdentifiers(t string, ids []string) Identifiers {
	identifiers := []Identifier{}

	for _, id := range ids {
		identifiers = append(identifiers, Identifier{
			Type: t,
			ID:   id,
		})
	}

	return identifiers
}

// Identifiers represents a slice of Identifier.
type Identifiers []Identifier

// IDs returns the IDs part of the Identifiers.
func (i Identifiers) IDs() []string {
	ids := make([]string, len(i))

	for n := range i {
		ids[n] = i[n].ID
	}

	return ids
}

// Identifier represents a resource's type and ID.
type Identifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// UnmarshalIdentifier reads a payload where the main data is one identifier to
// build and return an Identifier object.
//
// schema must not be nil.
func UnmarshalIdentifier(payload []byte, schema *Schema) (Identifier, error) {
	iden := Identifier{}

	err := json.Unmarshal(payload, &iden)
	if err != nil {
		return Identifier{}, err
	}

	switch {
	case iden.ID == "":
		return Identifier{}, errors.New("identifier has no ID")
	case iden.Type == "":
		return Identifier{}, errors.New("identifier has no type")
	case schema != nil && !schema.HasType(iden.Type):
		return Identifier{}, fmt.Errorf("type %q is unknown", iden.Type)
	}

	return iden, nil
}

// UnmarshalIdentifiers reads a payload where the main data is a collection of
// identifiers to build and return an Idenfitiers slice.
//
// schema must not be nil.
func UnmarshalIdentifiers(payload []byte, schema *Schema) (Identifiers, error) {
	raw := []*json.RawMessage{}

	err := json.Unmarshal(payload, &raw)
	if err != nil {
		return Identifiers{}, err
	}

	idens := make([]Identifier, len(raw))

	for i, r := range raw {
		iden, err := UnmarshalIdentifier(*r, schema)
		if err != nil {
			return nil, err
		}

		idens[i] = iden
	}

	return idens, nil
}
