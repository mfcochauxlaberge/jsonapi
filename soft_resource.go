package jsonapi

import (
	"fmt"
	"time"
)

// SoftResource represents a resource whose type is defined by an internal field
// of type *Type.
//
// Changing the type automatically changes the resource's attributes and
// relationships. When a field is added, its value is the zero value of the
// field's type.
type SoftResource struct {
	Type *Type

	id   string
	data map[string]interface{}
}

// Attrs returns the resource's attributes.
func (sr *SoftResource) Attrs() map[string]Attr {
	sr.check()

	return sr.Type.Attrs
}

// Rels returns the resource's relationships.
func (sr *SoftResource) Rels() map[string]Rel {
	sr.check()

	return sr.Type.Rels
}

// AddAttr adds an attribute.
func (sr *SoftResource) AddAttr(attr Attr) {
	sr.check()

	for _, name := range sr.fields() {
		if name == attr.Name {
			return
		}
	}

	sr.Type.Attrs[attr.Name] = attr
}

// AddRel adds a relationship.
func (sr *SoftResource) AddRel(rel Rel) {
	sr.check()

	for _, name := range sr.fields() {
		if name == rel.FromName {
			return
		}
	}

	sr.Type.Rels[rel.FromName] = rel
}

// RemoveField removes a field.
func (sr *SoftResource) RemoveField(field string) {
	sr.check()
	delete(sr.Type.Attrs, field)
	delete(sr.Type.Rels, field)
}

// Attr returns the attribute named after key.
func (sr *SoftResource) Attr(key string) Attr {
	sr.check()

	return sr.Type.Attrs[key]
}

// Rel returns the relationship named after key.
func (sr *SoftResource) Rel(key string) Rel {
	sr.check()

	return sr.Type.Rels[key]
}

// New returns a new resource (of type SoftResource) with the same type but
// without the values.
func (sr *SoftResource) New() Resource {
	sr.check()

	typ := sr.Type.Copy()

	return &SoftResource{
		Type: &typ,
	}
}

// GetID returns the resource's ID.
func (sr *SoftResource) GetID() string {
	sr.check()

	return sr.id
}

// GetType returns the resource's type.
func (sr *SoftResource) GetType() Type {
	sr.check()

	return *sr.Type
}

// Get returns the value associated to the field named after key.
func (sr *SoftResource) Get(key string) interface{} {
	sr.check()

	if _, ok := sr.Type.Attrs[key]; ok {
		if v, ok := sr.data[key]; ok {
			return v
		}
	}

	return nil
}

// SetID sets the resource's ID.
func (sr *SoftResource) SetID(id string) {
	sr.check()
	sr.id = id
}

// SetType sets the resource's type.
func (sr *SoftResource) SetType(typ *Type) {
	sr.check()
	sr.Type = typ
}

// Set sets the value associated to the field named key to v.
func (sr *SoftResource) Set(key string, v interface{}) {
	sr.check()

	if attr, ok := sr.Type.Attrs[key]; ok {
		if GetAttrTypeString(attr.Type, attr.Nullable) == fmt.Sprintf("%T", v) {
			sr.data[key] = v
		} else if v == nil && attr.Nullable {
			sr.data[key] = GetZeroValue(attr.Type, attr.Nullable)
		}
	}
}

// GetToOne returns the value associated to the relationship named after key.
func (sr *SoftResource) GetToOne(key string) string {
	sr.check()

	if _, ok := sr.Type.Rels[key]; ok {
		return sr.data[key].(string)
	}

	return ""
}

// GetToMany returns the value associated to the relationship named after key.
func (sr *SoftResource) GetToMany(key string) []string {
	sr.check()

	if _, ok := sr.Type.Rels[key]; ok {
		return sr.data[key].([]string)
	}

	return []string{}
}

// SetToOne sets the relationship named after key to rel.
func (sr *SoftResource) SetToOne(key string, v string) {
	sr.check()

	if rel, ok := sr.Type.Rels[key]; ok && rel.ToOne {
		sr.data[key] = v
	}
}

// SetToMany sets the relationship named after key to rel.
func (sr *SoftResource) SetToMany(key string, v []string) {
	sr.check()

	if rel, ok := sr.Type.Rels[key]; ok && !rel.ToOne {
		sr.data[key] = v
	}
}

// Copy return a new SoftResource object with the same type and values.
func (sr *SoftResource) Copy() Resource {
	sr.check()

	typ := sr.Type.Copy()

	return &SoftResource{
		Type: &typ,
		id:   sr.id,
		data: copyData(sr.data),
	}
}

func (sr *SoftResource) fields() []string {
	fields := make([]string, 0, len(sr.Type.Attrs)+len(sr.Type.Rels))
	for i := range sr.Type.Attrs {
		fields = append(fields, sr.Type.Attrs[i].Name)
	}

	for i := range sr.Type.Rels {
		fields = append(fields, sr.Type.Rels[i].FromName)
	}

	return fields
}

func (sr *SoftResource) check() {
	if sr.Type == nil {
		sr.Type = &Type{}
	}

	if sr.Type.Attrs == nil {
		sr.Type.Attrs = map[string]Attr{}
	}

	if sr.Type.Rels == nil {
		sr.Type.Rels = map[string]Rel{}
	}

	if sr.data == nil {
		sr.data = map[string]interface{}{}
	}

	for i := range sr.Type.Attrs {
		n := sr.Type.Attrs[i].Name
		if _, ok := sr.data[n]; !ok {
			sr.data[n] = GetZeroValue(sr.Type.Attrs[i].Type, sr.Type.Attrs[i].Nullable)
		}
	}

	for i := range sr.Type.Rels {
		n := sr.Type.Rels[i].FromName
		if _, ok := sr.data[n]; !ok {
			if sr.Type.Rels[i].ToOne {
				sr.data[n] = ""
			} else {
				sr.data[n] = []string{}
			}
		}
	}

	fields := sr.fields()
	if len(fields) < len(sr.data) {
		for k := range sr.data {
			found := false

			for _, f := range fields {
				if k == f {
					found = true
					break
				}
			}

			if !found {
				delete(sr.data, k)
			}
		}
	}
}

func copyData(d map[string]interface{}) map[string]interface{} {
	d2 := map[string]interface{}{}

	for k, v := range d {
		switch v2 := v.(type) {
		case string:
			d2[k] = v2
		case int:
			d2[k] = v2
		case int8:
			d2[k] = v2
		case int16:
			d2[k] = v2
		case int32:
			d2[k] = v2
		case int64:
			d2[k] = v2
		case uint:
			d2[k] = v2
		case uint8:
			d2[k] = v2
		case uint16:
			d2[k] = v2
		case uint32:
			d2[k] = v2
		case uint64:
			d2[k] = v2
		case bool:
			d2[k] = v2
		case time.Time:
			d2[k] = v2
		case []uint8:
			nv := make([]byte, len(v2))
			_ = copy(nv, v2)
			d2[k] = v2
		case []string:
			nv := make([]string, len(v2))
			_ = copy(nv, v2)
			d2[k] = v2
		case *string:
			d2[k] = v2
		case *int:
			d2[k] = v2
		case *int8:
			d2[k] = v2
		case *int16:
			d2[k] = v2
		case *int32:
			d2[k] = v2
		case *int64:
			d2[k] = v2
		case *uint:
			d2[k] = v2
		case *uint8:
			d2[k] = v2
		case *uint16:
			d2[k] = v2
		case *uint32:
			d2[k] = v2
		case *uint64:
			d2[k] = v2
		case *bool:
			d2[k] = v2
		case *time.Time:
			d2[k] = v2
		case *[]uint8:
			if v2 == nil {
				d2[k] = (*[]uint8)(nil)
			} else {
				nv := make([]byte, len(*v2))
				_ = copy(nv, *v2)
				d2[k] = v2
			}
		}
	}

	return d2
}
