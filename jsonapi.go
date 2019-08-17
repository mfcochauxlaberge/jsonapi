package jsonapi

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Marshal marshals a document according to the JSON:API speficication.
//
// Both doc and url must not be nil.
func Marshal(doc *Document, url *URL) ([]byte, error) {
	// Data
	var (
		data   json.RawMessage
		errors json.RawMessage
		err    error
	)

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
			url.Params.Fields[col.GetType().Name],
			doc.RelData,
		)
	} else if id, ok := doc.Data.(Identifier); ok {
		// Identifier
		data, err = json.Marshal(id)
	} else if ids, ok := doc.Data.(Identifiers); ok {
		// Identifiers
		data, err = json.Marshal(ids)
	} else if e, ok := doc.Data.(Error); ok {
		// Error
		errors, err = json.Marshal([]Error{e})
	} else if es, ok := doc.Data.([]Error); ok {
		// Errors
		errors, err = json.Marshal(es)
	} else {
		data = []byte("null")
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

// Unmarshal reads a payload to build and return a document object.
//
// Both url and schema must not be nil.
func Unmarshal(payload []byte, url *URL, schema *Schema) (*Document, error) {
	doc := &Document{}
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
		}
	}

	// Identifiers
	if url.RelKind == "self" {
		if !url.IsCol {
			inc := Identifier{}
			err = json.Unmarshal(ske.Data, &inc)
			if err != nil {
				return nil, err
			}
			doc.Data = inc
		} else {
			incs := Identifiers{}
			err = json.Unmarshal(ske.Data, &incs)
			if err != nil {
				return nil, err
			}
			doc.Data = incs
		}
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

		for i, inc2 := range incs {
			typ := schema.GetType(inc2.Type)
			res2 := &SoftResource{Type: &typ}
			err = json.Unmarshal(ske.Included[i], res2)
			if err != nil {
				return nil, err
			}
			doc.Included = append(doc.Included, res2)
		}
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
			if field == rel.Name {
				include = true
				break
			}
		}

		if include {
			var raw json.RawMessage

			if rel.ToOne {
				s := map[string]map[string]string{
					"links": buildRelationshipLinks(r, prepath, rel.Name),
				}

				for _, n := range relData[r.GetType().Name] {
					if n == rel.Name {
						id := r.GetToOne(rel.Name)
						if id != "" {
							s["data"] = map[string]string{
								"id":   r.GetToOne(rel.Name),
								"type": rel.Type,
							}
						} else {
							s["data"] = nil
						}
						break
					}
				}

				raw, _ = json.Marshal(s)
				rels[rel.Name] = &raw
			} else {
				s := map[string]interface{}{
					"links": buildRelationshipLinks(r, prepath, rel.Name),
				}

				for _, n := range relData[r.GetType().Name] {
					if n == rel.Name {
						data := []map[string]string{}
						ids := r.GetToMany(rel.Name)
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
				rels[rel.Name] = &raw
			}
		}
	}
	mapPl["relationships"] = rels

	// Links
	mapPl["links"] = map[string]string{
		"self": buildSelfLink(r, prepath), // TODO
	}

	pl, err := json.Marshal(mapPl)
	if err != nil {
		panic(fmt.Errorf("jsonapi: could not marshal resource: %s", err.Error()))
	}
	return pl
}

// marshalCollection marshals a Collection into a JSON-encoded payload.
func marshalCollection(c Collection, prepath string, fields []string, relData map[string][]string) []byte {
	var raws []*json.RawMessage

	if c.Len() == 0 {
		return []byte("[]")
	}

	for i := 0; i < c.Len(); i++ {
		r := c.At(i)
		raw := json.RawMessage(marshalResource(r, prepath, fields, relData))
		raws = append(raws, &raw)
	}

	pl, err := json.Marshal(raws)
	if err != nil {
		panic(fmt.Errorf("jsonapi: could not marshal collection: %s", err.Error()))
	}
	return pl
}

// unmarshalResource unmarshals a JSON-encoded payload into a Resource.
func unmarshalResource(data []byte, schema *Schema) (Resource, error) {
	var rske resourceSkeleton
	err := json.Unmarshal(data, &rske)
	if err != nil {
		return nil, err
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
					var val string
					err = json.Unmarshal(v.Data, &val)
					res.SetToOne(rel.Name, val)
				} else {
					var val []string
					err = json.Unmarshal(v.Data, &val)
					res.SetToMany(rel.Name, val)
				}
			}
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return nil, NewErrInvalidFieldValueInBody(
					rel.Name,
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
