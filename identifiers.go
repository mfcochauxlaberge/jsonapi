package jsonapi

import "encoding/json"

// NewIdentifiers ...
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

// Identifiers ...
type Identifiers []Identifier

// IDs ...
func (i Identifiers) IDs() []string {
	ids := []string{}

	for _, id := range i {
		ids = append(ids, id.ID)
	}

	return ids
}

// Identifier ...
type Identifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// MarshalIdentifiers ...
func MarshalIdentifiers(ids []string, toOne bool) json.RawMessage {
	raw := ""

	return []byte(raw)
}
