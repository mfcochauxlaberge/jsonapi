package jsonapi

import "encoding/json"

// marshalResource ...
func marshalResource(r Resource, host string, fields []string, relData map[string][]string) ([]byte, error) {
	mapPl := map[string]interface{}{}

	// ID and type
	mapPl["id"], mapPl["type"] = r.IDAndType()

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
					"links": buildRelationshipLinks(r, host, rel.Name),
				}

				for n, _ := range relData {
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
					"links": buildRelationshipLinks(r, host, rel.Name),
				}

				for n, _ := range relData {
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
		"self": buildSelfLink(r, host), // TODO
	}

	return json.Marshal(mapPl)
}

// marshalCollection ...
func marshalCollection(c Collection, host string, fields []string, relData map[string][]string) ([]byte, error) {
	var raws []*json.RawMessage

	if c.Len() == 0 {
		return []byte("[]"), nil
	}

	for i := 0; i < c.Len(); i++ {
		r := c.Elem(i)
		var raw json.RawMessage
		raw, err := marshalResource(r, host, fields, relData)
		if err != nil {
			return []byte{}, err
		}
		raws = append(raws, &raw)
	}

	return json.Marshal(raws)
}
