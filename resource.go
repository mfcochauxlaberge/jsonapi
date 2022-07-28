package jsonapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

// A Resource is an element of a collection.
type Resource interface {
	// Structure
	Attrs() map[string]Attr
	Rels() map[string]Rel

	// Read
	GetType() Type
	Get(key string) interface{}

	// Update
	Set(key string, val interface{})
}

// MarshalResource marshals a Resource into a JSON-encoded payload.
func MarshalResource(r Resource, prepath string, fields []string, relData map[string][]string) []byte {
	mapPl := map[string]interface{}{}

	mapPl["id"] = r.Get("id").(string)
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

	if len(attrs) > 0 {
		mapPl["attributes"] = attrs
	}

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
						id := r.Get(rel.FromName).(string)
						if id != "" {
							s["data"] = map[string]string{
								"id":   r.Get(rel.FromName).(string),
								"type": rel.ToType,
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
						ids := r.Get(rel.FromName).([]string)
						sort.Strings(ids)
						for _, id := range ids {
							data = append(data, map[string]string{
								"id":   id,
								"type": rel.ToType,
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

	if len(rels) > 0 {
		mapPl["relationships"] = rels
	}

	// Links
	mapPl["links"] = map[string]string{
		"self": buildSelfLink(r, prepath),
	}

	// Meta
	if m, ok := r.(MetaHolder); ok {
		if len(m.Meta()) > 0 {
			mapPl["meta"] = m.Meta()
		}
	}

	// NOTE An error should not happen.
	pl, _ := json.Marshal(mapPl)

	return pl
}

// UnmarshalResource unmarshals a JSON-encoded payload into a Resource.
func UnmarshalResource(data []byte, schema *Schema) (Resource, error) {
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

	res.Set("id", rske.ID)

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
					var iden Identifier
					err = json.Unmarshal(v.Data, &iden)
					res.Set(rel.FromName, iden.ID)
				} else {
					var idens Identifiers
					err = json.Unmarshal(v.Data, &idens)
					ids := make([]string, len(idens))
					for i := range idens {
						ids[i] = idens[i].ID
					}
					res.Set(rel.FromName, ids)
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

	// Meta
	if m, ok := res.(MetaHolder); ok {
		m.SetMeta(rske.Meta)
	}

	return res, nil
}

// UnmarshalPartialResource unmarshals the given payload into a *SoftResource.
//
// The returned *SoftResource will only contain the information found in the
// payload. That means that fields not in the payload won't be part of the
// *SoftResource. Its type will be a new type whose fields will be a subset of
// the fields of the corresponding type from the schema.
//
// This is useful when handling a PATCH request where only some fields might be
// set to a value. UnmarshalResource returns a Resource where the missing fields
// are added and set to their zero value, but UnmarshalPartialResource does not
// do that. Therefore, the user is able to tell which fields have been set.
func UnmarshalPartialResource(data []byte, schema *Schema) (*SoftResource, error) {
	var rske resourceSkeleton
	err := json.Unmarshal(data, &rske)

	if err != nil {
		return nil, NewErrBadRequest(
			"Invalid JSON",
			"The provided JSON body could not be read.",
		)
	}

	typ := schema.GetType(rske.Type)
	newType := Type{
		Name: typ.Name,
	}
	res := &SoftResource{
		Type: &newType,
		id:   rske.ID,
	}

	for a, v := range rske.Attributes {
		if attr, ok := typ.Attrs[a]; ok {
			val, err := attr.UnmarshalToType(v)
			if err != nil {
				return nil, err
			}

			_ = newType.AddAttr(attr)
			res.Set(attr.Name, val)
		} else {
			return nil, NewErrUnknownFieldInBody(typ.Name, a)
		}
	}

	for r, v := range rske.Relationships {
		if rel, ok := typ.Rels[r]; ok {
			if len(v.Data) > 0 {
				if rel.ToOne {
					var iden Identifier
					err = json.Unmarshal(v.Data, &iden)
					_ = newType.AddRel(rel)
					res.Set(rel.FromName, iden.ID)
				} else {
					var idens Identifiers
					err = json.Unmarshal(v.Data, &idens)
					ids := make([]string, len(idens))
					for i := range idens {
						ids[i] = idens[i].ID
					}
					_ = newType.AddRel(rel)
					res.Set(rel.FromName, ids)
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

// Equal reports whether r1 and r2 are equal.
//
// Two resources are equal if their types are equal, all the attributes are
// equal (same type and same value), and all the relationstips are equal.
//
// IDs are ignored.
func Equal(r1, r2 Resource) bool {
	// Type
	if r1.GetType().Name != r2.GetType().Name {
		return false
	}

	// Attributes
	attrs := r1.Attrs()
	r1Attrs := make([]Attr, 0, len(attrs))

	for name := range attrs {
		r1Attrs = append(r1Attrs, attrs[name])
	}

	sort.Slice(r1Attrs, func(i, j int) bool {
		return r1Attrs[i].Name < r1Attrs[j].Name
	})

	attrs = r2.Attrs()
	r2Attrs := make([]Attr, 0, len(attrs))

	for name := range attrs {
		r2Attrs = append(r2Attrs, attrs[name])
	}

	sort.Slice(r2Attrs, func(i, j int) bool {
		return r2Attrs[i].Name < r2Attrs[j].Name
	})

	if len(r1Attrs) != len(r2Attrs) {
		return false
	}

	for i, attr1 := range r1Attrs {
		attr2 := r2Attrs[i]
		if !reflect.DeepEqual(r1.Get(attr1.Name), r2.Get(attr2.Name)) {
			// TODO Fix the following condition one day. Basically, all
			// nils (nil pointer, nil slice, etc) should be considered
			// equal to a nil empty interface.
			if fmt.Sprintf("%v", r1.Get(attr1.Name)) == "<nil>" &&
				fmt.Sprintf("%v", r2.Get(attr1.Name)) == "<nil>" {
				continue
			}

			return false
		}
	}

	// Relationships
	rels := r1.Rels()
	r1Rels := make([]Rel, 0, len(rels))

	for name := range rels {
		r1Rels = append(r1Rels, rels[name])
	}

	sort.Slice(r1Rels, func(i, j int) bool {
		return r1Rels[i].FromName < r1Rels[j].FromName
	})

	rels = r2.Rels()
	r2Rels := make([]Rel, 0, len(rels))

	for name := range rels {
		r2Rels = append(r2Rels, rels[name])
	}

	sort.Slice(r2Rels, func(i, j int) bool {
		return r2Rels[i].FromName < r2Rels[j].FromName
	})

	if len(r1Rels) != len(r2Rels) {
		return false
	}

	for i, rel1 := range r1Rels {
		rel2 := r2Rels[i]
		if rel1.ToOne != rel2.ToOne {
			return false
		}

		if rel1.ToOne {
			if r1.Get(rel1.FromName).(string) != r2.Get(rel2.FromName).(string) {
				return false
			}
		} else {
			v1 := r1.Get(rel1.FromName).([]string)
			v2 := r2.Get(rel2.FromName).([]string)
			if len(v1) != 0 || len(v2) != 0 {
				if !reflect.DeepEqual(v1, v2) {
					return false
				}
			}
		}
	}

	return true
}

// EqualStrict is like Equal, but it also considers IDs.
func EqualStrict(r1, r2 Resource) bool {
	if r1.Get("id").(string) != r2.Get("id").(string) {
		return false
	}

	return Equal(r1, r2)
}
