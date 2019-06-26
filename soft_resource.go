package jsonapi

import (
	"fmt"

	"github.com/mitchellh/copystructure"
)

var _ Resource = (*SoftResource)(nil)

// NewSoftResource returns a new SoftResource with the given type.
//
// It is also populated with values from vals.
func NewSoftResource(typ Type, vals map[string]interface{}) *SoftResource {
	res := &SoftResource{}
	res.typ = &typ

	for _, attr := range typ.Attrs {
		if val, ok := vals[attr.Name]; ok {
			res.Set(attr.Name, val)
		}
	}
	for _, rel := range typ.Rels {
		if val, ok := vals[rel.Name]; ok {
			res.Set(rel.Name, val)
		}
	}

	return res
}

// SoftResource represents a resource whose type is defined by an internal
// field of type *Type.
//
// Changing the type automatically changes the resource's attributes and
// relationships. When a field is added, its value is the zero value of the
// field's type.
type SoftResource struct {
	id   string
	typ  *Type
	data map[string]interface{}
}

// Attrs returns the resource's attributes.
func (sr *SoftResource) Attrs() map[string]Attr {
	sr.check()
	return sr.typ.Attrs
}

// Rels returns the resource's relationships.
func (sr *SoftResource) Rels() map[string]Rel {
	sr.check()
	return sr.typ.Rels
}

// AddAttr adds an attribute.
func (sr *SoftResource) AddAttr(attr Attr) {
	sr.check()
	for _, name := range sr.fields() {
		if name == attr.Name {
			return
		}
	}
	sr.typ.Attrs[attr.Name] = attr
}

// AddRel adds a relationship.
func (sr *SoftResource) AddRel(rel Rel) {
	sr.check()
	for _, name := range sr.fields() {
		if name == rel.Name {
			return
		}
	}
	sr.typ.Rels[rel.Name] = rel
}

// RemoveField removes a field.
func (sr *SoftResource) RemoveField(field string) {
	sr.check()
	delete(sr.typ.Attrs, field)
	delete(sr.typ.Rels, field)
}

// Attr returns the attribute named after key.
func (sr *SoftResource) Attr(key string) Attr {
	sr.check()
	return sr.typ.Attrs[key]
}

// Rel returns the relationship named after key.
func (sr *SoftResource) Rel(key string) Rel {
	sr.check()
	return sr.typ.Rels[key]
}

// New returns a new resource (of type SoftResource) with the same type
// but without the values.
func (sr *SoftResource) New() Resource {
	sr.check()
	return &SoftResource{
		typ: copystructure.Must(copystructure.Copy(sr.typ)).(*Type),
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
	return *sr.typ
}

// Get returns the value associated to the field named after key.
func (sr *SoftResource) Get(key string) interface{} {
	sr.check()
	if attr, ok := sr.typ.Attrs[key]; ok {
		if v, ok := sr.data[key]; ok {
			return v
		}
		return GetZeroValue(attr.Type, attr.Null)
	}
	if rel, ok := sr.typ.Rels[key]; ok {
		if v, ok := sr.data[key]; ok {
			return v
		}
		if rel.ToOne {
			return ""
		}
		return []string{}
	}
	return nil
}

// SetID sets the resource's ID.
func (sr *SoftResource) SetID(id string) {
	sr.check()
	sr.id = id
}

// SetType ...
func (sr *SoftResource) SetType(typ *Type) {
	sr.check()
	sr.typ = typ
}

// Set sets the value associated to the field named key to v.
func (sr *SoftResource) Set(key string, v interface{}) {
	sr.check()
	fmt.Printf("about to set %s to %v (%T)\n", key, v, v)
	if attr, ok := sr.typ.Attrs[key]; ok {
		fmt.Printf("attr found, type=%s\n", GetAttrTypeString(attr.Type, attr.Null))
		if GetAttrTypeString(attr.Type, attr.Null) == fmt.Sprintf("%T", v) {
			fmt.Printf("done (1)\n")
			sr.data[key] = v
		} else if v == nil && attr.Null {
			fmt.Printf("done (2)\n")
			sr.data[key] = GetZeroValue(attr.Type, attr.Null)
		}
	}
}

// GetToOne returns the value associated to the relationship named after key.
func (sr *SoftResource) GetToOne(key string) string {
	sr.check()
	if _, ok := sr.typ.Rels[key]; ok {
		return sr.data[key].(string)
	}
	return ""
}

// GetToMany returns the value associated to the relationship named after key.
func (sr *SoftResource) GetToMany(key string) []string {
	sr.check()
	if _, ok := sr.typ.Rels[key]; ok {
		return sr.data[key].([]string)
	}
	return []string{}
}

// SetToOne sets the relationship named after key to rel.
func (sr *SoftResource) SetToOne(key string, v string) {
	sr.check()
	if rel, ok := sr.typ.Rels[key]; ok && rel.ToOne {
		sr.data[key] = v
	}
}

// SetToMany sets the relationship named after key to rel.
func (sr *SoftResource) SetToMany(key string, v []string) {
	sr.check()
	if rel, ok := sr.typ.Rels[key]; ok && !rel.ToOne {
		sr.data[key] = v
	}
}

// Validate returns validation errors found in the resource.
func (sr *SoftResource) Validate() []error {
	sr.check()
	return []error{}
}

// Copy return a new SoftResource object with the same type and values.
func (sr *SoftResource) Copy() Resource {
	sr.check()
	return &SoftResource{
		id:   sr.id,
		typ:  copystructure.Must(copystructure.Copy(sr.typ)).(*Type),
		data: copystructure.Must(copystructure.Copy(sr.data)).(map[string]interface{}),
	}
}

// UnmarshalJSON parses the payload and populates a SoftResource.
func (sr *SoftResource) UnmarshalJSON(payload []byte) error {
	sr.check()
	// TODO
	return nil
}

func (sr *SoftResource) fields() []string {
	fields := make([]string, 0, len(sr.typ.Attrs)+len(sr.typ.Rels))
	for i := range sr.typ.Attrs {
		fields = append(fields, sr.typ.Attrs[i].Name)
	}
	for i := range sr.typ.Rels {
		fields = append(fields, sr.typ.Rels[i].Name)
	}
	return fields
}

func (sr *SoftResource) check() {
	if sr.typ == nil {
		sr.typ = &Type{}
	}
	if sr.typ.Attrs == nil {
		sr.typ.Attrs = map[string]Attr{}
	}
	if sr.typ.Rels == nil {
		sr.typ.Rels = map[string]Rel{}
	}
	if sr.data == nil {
		sr.data = map[string]interface{}{}
	}

	for i := range sr.typ.Attrs {
		n := sr.typ.Attrs[i].Name
		if _, ok := sr.data[n]; !ok {
			sr.data[n] = GetZeroValue(sr.typ.Attrs[i].Type, sr.typ.Attrs[i].Null)
		}
	}
	for i := range sr.typ.Rels {
		n := sr.typ.Rels[i].Name
		if _, ok := sr.data[n]; !ok {
			if sr.typ.Rels[i].ToOne {
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
