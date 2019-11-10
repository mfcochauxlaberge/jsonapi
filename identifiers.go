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
	ids := []string{}

	for _, id := range i {
		ids = append(ids, id.ID)
	}

	return ids
}

// Identifier represents a resource's type and ID.
type Identifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// UnmarshalIdentifiers reads a payload where the main data is one or more
// identifiers to build and return a Document object.
//
// The included top-level member is ignored.
//
// schema must not be nil.
func UnmarshalIdentifiers(payload []byte, schema *Schema) (*Document, error) {
	doc := &Document{
		Included:  []Resource{},
		Resources: map[string]map[string]struct{}{},
		Links:     map[string]Link{},
		RelData:   map[string][]string{},
		Meta:      map[string]interface{}{},
	}
	ske := &payloadSkeleton{}

	// Unmarshal
	err := json.Unmarshal(payload, ske)
	if err != nil {
		return nil, err
	}

	// Identifiers
	if len(ske.Data) > 0 {
		if ske.Data[0] == '{' {
			inc := Identifier{}
			err = json.Unmarshal(ske.Data, &inc)
			if err != nil {
				return nil, err
			}
			doc.Data = inc
		} else if ske.Data[0] == '[' {
			incs := Identifiers{}
			err = json.Unmarshal(ske.Data, &incs)
			if err != nil {
				return nil, err
			}
			doc.Data = incs
		}
	} else if len(ske.Errors) > 0 {
		doc.Errors = ske.Errors
	} else {
		return nil, NewErrMissingDataMember()
	}

	// Meta
	doc.Meta = ske.Meta

	return doc, nil
}
