package jsonapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// marshalResource ...
func marshalResource(r Resource, prepath string, fields []string, relData map[string][]string) ([]byte, error) {
	mapPl := map[string]interface{}{}

	// ID and type
	mapPl["id"] = r.GetID()
	mapPl["type"] = r.GetType()

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

	return json.Marshal(mapPl)
}

// marshalCollection ...
func marshalCollection(c Collection, prepath string, fields []string, relData map[string][]string) ([]byte, error) {
	var raws []*json.RawMessage

	if c.Len() == 0 {
		return []byte("[]"), nil
	}

	for i := 0; i < c.Len(); i++ {
		r := c.Elem(i)
		var raw json.RawMessage
		raw, err := marshalResource(r, prepath, fields, relData)
		if err != nil {
			return []byte{}, err
		}
		raws = append(raws, &raw)
	}

	return json.Marshal(raws)
}

// ReflectType takes a struct or a pointer to a struct to analyse and
// builds a Type object that is returned.
//
// If an error is returned, the Type object will be empty.
func ReflectType(v interface{}) (Type, error) {
	typ := Type{}

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return typ, errors.New("jsonapi: value must represent a struct")
	}

	err := CheckType(val.Interface())
	if err != nil {
		return typ, fmt.Errorf("jsonapi: invalid type: %s", err)
	}

	// ID and type
	_, typ.Name = IDAndType(v)

	// Attributes
	typ.Attrs = map[string]Attr{}
	for i := 0; i < val.NumField(); i++ {
		fs := val.Type().Field(i)
		jsonTag := fs.Tag.Get("json")
		apiTag := fs.Tag.Get("api")

		if apiTag == "attr" {
			fieldType, null := GetAttrType(fs.Type.String())
			typ.Attrs[jsonTag] = Attr{
				Name: jsonTag,
				Type: fieldType,
				Null: null,
			}
		}
	}

	// Relationships
	typ.Rels = map[string]Rel{}
	for i := 0; i < val.NumField(); i++ {
		fs := val.Type().Field(i)
		jsonTag := fs.Tag.Get("json")
		relTag := strings.Split(fs.Tag.Get("api"), ",")
		invName := ""
		if len(relTag) == 3 {
			invName = relTag[2]
		}

		toOne := true
		if fs.Type.String() == "[]string" {
			toOne = false
		}

		if relTag[0] == "rel" {
			typ.Rels[jsonTag] = Rel{
				Name:        jsonTag,
				Type:        relTag[1],
				ToOne:       toOne,
				InverseName: invName,
				InverseType: typ.Name,
			}
		}
	}

	return typ, nil
}
