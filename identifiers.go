package jsonapi

import "encoding/json"

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
	// TODO Validate with schema.
	return iden, err
}

// UnmarshalIdentifiers reads a payload where the main data is a collection of
// identifiers to build and return an Idenfitiers slice.
//
// schema must not be nil.
func UnmarshalIdentifiers(payload []byte, schema *Schema) (Identifiers, error) {
	idens := Identifiers{}
	err := json.Unmarshal(payload, &idens)
	// TODO Validate with schema.
	return idens, err
}
