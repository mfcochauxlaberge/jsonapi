package jsonapi

import (
	"encoding/json"
	"errors"
	"sort"
)

// Marshal marshals a document according to the JSON:API speficication.
//
// Both doc and url must not be nil.
func Marshal(doc *Document, url *URL) ([]byte, error) {
	var err error

	// Data
	var data json.RawMessage
	if res, ok := doc.Data.(Resource); ok {
		// Resource
		data = marshalResource(
			res,
			doc.PrePath,
			url.Params.Fields[res.GetType().Name],
			doc.RelData,
		)
	} else if col, ok := doc.Data.(Collection); ok {
		// Collection
		data = marshalCollection(
			col,
			doc.PrePath,
			url.Params.Fields,
			doc.RelData,
		)
	} else if id, ok := doc.Data.(Identifier); ok {
		// Identifier
		data, err = json.Marshal(id)
	} else if ids, ok := doc.Data.(Identifiers); ok {
		// Identifiers
		data, err = json.Marshal(ids)
	} else if doc.Data != nil {
		err = errors.New("data contains an unknown type")
	} else if len(doc.Errors) == 0 {
		data = []byte("null")
	}

	// Data
	var errors json.RawMessage
	if len(doc.Errors) > 0 {
		// Errors
		errors, err = json.Marshal(doc.Errors)
	}

	if err != nil {
		return []byte{}, err
	}

	// Included
	var inclusions []*json.RawMessage
	if len(doc.Included) > 0 {
		sort.Slice(doc.Included, func(i, j int) bool {
			return doc.Included[i].GetID() < doc.Included[j].GetID()
		})

		if len(data) > 0 {
			for key := range doc.Included {
				typ := doc.Included[key].GetType().Name
				raw := marshalResource(
					doc.Included[key],
					doc.PrePath,
					url.Params.Fields[typ],
					doc.RelData,
				)
				rawm := json.RawMessage(raw)
				inclusions = append(inclusions, &rawm)
			}
		}
	}

	// Marshaling
	plMap := map[string]interface{}{}

	if len(errors) > 0 {
		plMap["errors"] = errors
	} else if len(data) > 0 {
		plMap["data"] = data

		if len(inclusions) > 0 {
			plMap["included"] = inclusions
		}
	}

	if len(doc.Meta) > 0 {
		plMap["meta"] = doc.Meta
	}

	if url != nil {
		plMap["links"] = map[string]string{
			"self": doc.PrePath + url.String(),
		}
	}
	plMap["jsonapi"] = map[string]string{"version": "1.0"}

	return json.Marshal(plMap)
}

// Unmarshal reads a payload to build and return a Document object.
//
// schema must not be nil.
func Unmarshal(payload []byte, schema *Schema) (*Document, error) {
	doc := NewDocument()
	ske := &payloadSkeleton{}

	// Unmarshal
	err := json.Unmarshal(payload, ske)
	if err != nil {
		return nil, err
	}

	// Data
	if len(ske.Data) > 0 {
		if ske.Data[0] == '{' {
			// Resource
			res, err := unmarshalResource(ske.Data, schema)
			if err != nil {
				return nil, err
			}
			doc.Data = res
		} else if ske.Data[0] == '[' {
			col, err := unmarshalCollection(ske.Data, schema)
			if err != nil {
				return nil, err
			}
			doc.Data = col
		} else if string(ske.Data) == "null" {
			doc.Data = nil
		} else {
			// TODO Not exactly the right error
			return nil, NewErrMissingDataMember()
		}
	} else if len(ske.Errors) > 0 {
		doc.Errors = ske.Errors
	} else {
		return nil, NewErrMissingDataMember()
	}

	// Included
	if len(ske.Included) > 0 {
		inc := Identifier{}
		incs := []Identifier{}
		for _, rawInc := range ske.Included {
			err = json.Unmarshal(rawInc, &inc)
			if err != nil {
				return nil, err
			}
			incs = append(incs, inc)
		}

		for i := range incs {
			res, err := unmarshalResource(ske.Included[i], schema)
			if err != nil {
				return nil, err
			}
			doc.Included = append(doc.Included, res)
		}
	}

	// Meta
	doc.Meta = ske.Meta

	return doc, nil
}

// UnmarshalIdentifiers reads a payload where the main data is one or more
// identifiers to build and return a Document object.
//
// The included top-level member is ignored.
//
// schema must not be nil.
func UnmarshalIdentifiers(payload []byte, schema *Schema) (*Document, error) {
	doc := NewDocument()
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

// marshalResource marshals a Resource into a JSON-encoded payload.
func marshalResource(r Resource, prepath string, fields []string, relData map[string][]string) []byte {
	mapPl := map[string]interface{}{}

	mapPl["id"] = r.GetID()
	mapPl["type"] = r.GetType().Name

	// Attributes
	attrs := map[string]interface{}{}
	for _, attr := range r.Attrs() {
		for _, field := range fields {
			if field == attr.Name {
				attrs[attr.Name] = r.Get(attr.Name)
				break
			}
		}
	}
	mapPl["attributes"] = attrs

	// Relationships
	rels := map[string]*json.RawMessage{}
	for _, rel := range r.Rels() {
		include := false
		for _, field := range fields {
			if field == rel.FromName {
				include = true
				break
			}
		}

		if include {
			var raw json.RawMessage

			if rel.ToOne {
				s := map[string]map[string]string{
					"links": buildRelationshipLinks(r, prepath, rel.FromName),
				}

				for _, n := range relData[r.GetType().Name] {
					if n == rel.FromName {
						id := r.GetToOne(rel.FromName)
						if id != "" {
							s["data"] = map[string]string{
								"id":   r.GetToOne(rel.FromName),
								"type": rel.Type,
							}
						} else {
							s["data"] = nil
						}
						break
					}
				}

				raw, _ = json.Marshal(s)
				rels[rel.FromName] = &raw
			} else {
				s := map[string]interface{}{
					"links": buildRelationshipLinks(r, prepath, rel.FromName),
				}

				for _, n := range relData[r.GetType().Name] {
					if n == rel.FromName {
						data := []map[string]string{}
						ids := r.GetToMany(rel.FromName)
						sort.Strings(ids)
						for _, id := range ids {
							data = append(data, map[string]string{
								"id":   id,
								"type": rel.Type,
							})
						}
						s["data"] = data
						break
					}
				}

				raw, _ = json.Marshal(s)
				rels[rel.FromName] = &raw
			}
		}
	}
	mapPl["relationships"] = rels

	// Links
	mapPl["links"] = map[string]string{
		"self": buildSelfLink(r, prepath), // TODO
	}

	// NOTE An error should not happen.
	pl, _ := json.Marshal(mapPl)
	return pl
}

// marshalCollection marshals a Collection into a JSON-encoded payload.
func marshalCollection(c Collection, prepath string, fields map[string][]string, relData map[string][]string) []byte {
	var raws []*json.RawMessage

	if c.Len() == 0 {
		return []byte("[]")
	}

	for i := 0; i < c.Len(); i++ {
		r := c.At(i)
		raw := json.RawMessage(
			marshalResource(r, prepath, fields[r.GetType().Name], relData),
		)
		raws = append(raws, &raw)
	}

	// NOTE An error should not happen.
	pl, _ := json.Marshal(raws)
	return pl
}

// unmarshalResource unmarshals a JSON-encoded payload into a Resource.
func unmarshalResource(data []byte, schema *Schema) (Resource, error) {
	var rske resourceSkeleton
	err := json.Unmarshal(data, &rske)
	if err != nil {
		return nil, NewErrBadRequest(
			"Invalid JSON",
			"The provided JSON body could not be read.",
		)
	}

	typ := schema.GetType(rske.Type)
	res := typ.New()

	res.SetID(rske.ID)

	for a, v := range rske.Attributes {
		if attr, ok := typ.Attrs[a]; ok {
			val, err := attr.UnmarshalToType(v)
			if err != nil {
				return nil, err
			}
			res.Set(attr.Name, val)
		} else {
			return nil, NewErrUnknownFieldInBody(typ.Name, a)
		}
	}
	for r, v := range rske.Relationships {
		if rel, ok := typ.Rels[r]; ok {
			if len(v.Data) > 0 {
				if rel.ToOne {
					var iden identifierSkeleton
					err = json.Unmarshal(v.Data, &iden)
					res.SetToOne(rel.FromName, iden.ID)
				} else {
					var idens []identifierSkeleton
					err = json.Unmarshal(v.Data, &idens)
					ids := make([]string, len(idens))
					for i := range idens {
						ids[i] = idens[i].ID
					}
					res.SetToMany(rel.FromName, ids)
				}
			}
			if err != nil {
				return nil, NewErrInvalidFieldValueInBody(
					rel.FromName,
					string(v.Data),
					typ.Name,
				)
			}
		} else {
			return nil, NewErrUnknownFieldInBody(typ.Name, r)
		}
	}

	return res, nil
}

// unmarshalCollection unmarshals a JSON-encoded payload into a Collection.
func unmarshalCollection(data []byte, schema *Schema) (Collection, error) {
	var cske []json.RawMessage
	err := json.Unmarshal(data, &cske)
	if err != nil {
		return nil, err
	}

	col := &Resources{}
	for i := range cske {
		res, err := unmarshalResource(cske[i], schema)
		if err != nil {
			return nil, err
		}
		col.Add(res)
	}

	return col, nil
}
