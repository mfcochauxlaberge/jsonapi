package jsonapi

import (
	"github.com/mitchellh/copystructure"
)

// NewSoftResource ...
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

// SoftResource ...
type SoftResource struct {
	id   string
	typ  *Type
	data map[string]interface{}
}

// Attrs ...
func (sr *SoftResource) Attrs() map[string]Attr {
	sr.check()
	return sr.typ.Attrs
}

// Rels ...
func (sr *SoftResource) Rels() map[string]Rel {
	sr.check()
	return sr.typ.Rels
}

// AddAttr ...
func (sr *SoftResource) AddAttr(attr Attr) {
	sr.check()
	for _, name := range sr.fields() {
		if name == attr.Name {
			return
		}
	}
	sr.typ.Attrs[attr.Name] = attr
}

// AddRel ...
func (sr *SoftResource) AddRel(rel Rel) {
	sr.check()
	for _, name := range sr.fields() {
		if name == rel.Name {
			return
		}
	}
	sr.typ.Rels[rel.Name] = rel
}

// RemoveField ...
func (sr *SoftResource) RemoveField(field string) {
	sr.check()
	delete(sr.typ.Attrs, field)
	delete(sr.typ.Rels, field)
}

// Attr ...
func (sr *SoftResource) Attr(key string) Attr {
	sr.check()
	return sr.typ.Attrs[key]
}

// Rel ...
func (sr *SoftResource) Rel(key string) Rel {
	sr.check()
	return sr.typ.Rels[key]
}

// New ...
func (sr *SoftResource) New() Resource {
	sr.check()
	return &SoftResource{
		typ: copystructure.Must(copystructure.Copy(sr.typ)).(*Type),
	}
}

// GetID ...
func (sr *SoftResource) GetID() string {
	sr.check()
	return sr.id
}

// GetType ...
func (sr *SoftResource) GetType() string {
	sr.check()
	return sr.typ.Name
}

// Get ...
func (sr *SoftResource) Get(key string) interface{} {
	sr.check()
	if _, ok := sr.typ.Attrs[key]; ok {
		return sr.data[key]
	}
	if _, ok := sr.typ.Rels[key]; ok {
		return sr.data[key]
	}
	return nil
}

// SetID ...
func (sr *SoftResource) SetID(id string) {
	sr.check()
	sr.id = id
}

// SetType ...
func (sr *SoftResource) SetType(typ *Type) {
	sr.check()
	sr.typ = typ
}

// Set ...
func (sr *SoftResource) Set(key string, v interface{}) {
	sr.check()
	if _, ok := sr.data[key]; ok {
		sr.data[key] = v
	}
}

// GetToOne ...
func (sr *SoftResource) GetToOne(key string) string {
	sr.check()
	if _, ok := sr.typ.Rels[key]; ok {
		return sr.data[key].(string)
	}
	return ""
}

// GetToMany ...
func (sr *SoftResource) GetToMany(key string) []string {
	sr.check()
	if _, ok := sr.typ.Rels[key]; ok {
		return sr.data[key].([]string)
	}
	return []string{}
}

// SetToOne ...
func (sr *SoftResource) SetToOne(key string, rel string) {
	sr.check()
	if _, ok := sr.data[key]; ok {
		sr.data[key] = rel
	}
}

// SetToMany ...
func (sr *SoftResource) SetToMany(key string, rels []string) {
	sr.check()
	if _, ok := sr.data[key]; ok {
		sr.data[key] = rels
	}
}

// Validate ...
func (sr *SoftResource) Validate() []error {
	sr.check()
	return []error{}
}

// Copy ...
func (sr *SoftResource) Copy() Resource {
	sr.check()
	return &SoftResource{
		id:   sr.id,
		typ:  copystructure.Must(copystructure.Copy(sr.typ)).(*Type),
		data: copystructure.Must(copystructure.Copy(sr.data)).(map[string]interface{}),
	}
}

// UnmarshalJSON ...
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
