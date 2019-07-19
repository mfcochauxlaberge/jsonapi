package jsonapi

import (
	"encoding/json"
	"fmt"
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
		data = marshalResource(res, doc.PrePath, url.Params.Fields[res.GetType().Name], doc.RelData)
	} else if col, ok := doc.Data.(Collection); ok {
		// Collection
		data = marshalCollection(col, doc.PrePath, url.Params.Fields[col.GetType().Name], doc.RelData)
	} else if id, ok := doc.Data.(Identifier); ok {
		// Identifer
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
	inclusions := []*json.RawMessage{}
	if len(data) > 0 {
		for key := range doc.Included {
			typ := doc.Included[key].GetType().Name
			raw := marshalResource(doc.Included[key], doc.PrePath, url.Params.Fields[typ], doc.RelData)
			rawm := json.RawMessage(raw)
			inclusions = append(inclusions, &rawm)
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
			"self": doc.PrePath + url.FullURL(),
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
	if !url.IsCol && url.RelKind == "" {
		typ := schema.GetType(url.ResType)
		res := &SoftResource{Type: &typ}
		err = json.Unmarshal(ske.Data, res)
		if err != nil {
			return nil, err
		}
		doc.Data = res
	} else if url.RelKind == "self" {
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

// marshalResource ...
func marshalResource(r Resource, prepath string, fields []string, relData map[string][]string) []byte {
	mapPl := map[string]interface{}{}

	// ID and type
	mapPl["id"] = r.GetID()
	mapPl["type"] = r.GetType().Name

	// Attributes
	attrs := map[string]interface{}{}
	for _, attr := range r.Attrs() {
		if len(fields) == 0 {
			attrs[attr.Name] = r.Get(attr.Name)
		} else {
			for _, field := range fields {
				if field == attr.Name {
					attrs[attr.Name] = r.Get(attr.Name)
					break
				}
			}
		}
	}
	mapPl["attributes"] = attrs

	// Relationships
	rels := map[string]*json.RawMessage{}
	for _, rel := range r.Rels() {
		include := false
		if len(fields) == 0 {
			include = true
		} else {
			for _, field := range fields {
				if field == rel.Name {
					include = true
					break
				}
			}
		}

		if include {
			if rel.ToOne {
				var raw json.RawMessage

				s := map[string]map[string]string{
					"links": buildRelationshipLinks(r, prepath, rel.Name),
				}

				for n := range relData {
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

				// var links map[string]string{}
				raw, _ = json.Marshal(s)
				rels[rel.Name] = &raw
			} else {
				var raw json.RawMessage

				s := map[string]interface{}{
					"links": buildRelationshipLinks(r, prepath, rel.Name),
				}

				for n := range relData {
					if n == rel.Name {
						data := []map[string]string{}

						for _, id := range r.GetToMany(rel.Name) {
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

// marshalCollection ...
func marshalCollection(c Collection, prepath string, fields []string, relData map[string][]string) []byte {
	var raws []*json.RawMessage

	if c.Len() == 0 {
		return []byte("[]")
	}

	for i := 0; i < c.Len(); i++ {
		r := c.At(i)
		var raw json.RawMessage
		raw = marshalResource(r, prepath, fields, relData)
		raws = append(raws, &raw)
	}

	pl, err := json.Marshal(raws)
	if err != nil {
		panic(fmt.Errorf("jsonapi: could not marshal collection: %s", err.Error()))
	}
	return pl
}
