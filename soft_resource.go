package jsonapi

import (
	"github.com/mitchellh/copystructure"
)

// SoftResource ...
type SoftResource struct {
	typ   string
	id    string
	attrs []Attr
	rels  []Rel
	data  map[string]interface{}
}

// Attrs ...
func (sr *SoftResource) Attrs() []Attr {
	sr.check()
	return sr.attrs
}

// Rels ...
func (sr *SoftResource) Rels() []Rel {
	sr.check()
	return sr.rels
}

// AddAttr ...
func (sr *SoftResource) AddAttr(attr Attr) {
	sr.check()
	for _, f := range sr.fields() {
		if f == attr.Name {
			return
		}
	}
	sr.attrs = append(sr.attrs, attr)
}

// AddRel ...
func (sr *SoftResource) AddRel(rel Rel) {
	sr.check()
	for _, f := range sr.fields() {
		if f == rel.Name {
			return
		}
	}
	sr.rels = append(sr.rels, rel)
}

// RemoveField ...
func (sr *SoftResource) RemoveField(field string) {
	sr.check()
	for i, a := range sr.attrs {
		if field == a.Name {
			sr.attrs = append(sr.attrs[:i], sr.attrs[i+1:]...)
			return
		}
	}
	for i, r := range sr.rels {
		if field == r.Name {
			sr.rels = append(sr.rels[:i], sr.rels[i+1:]...)
			return
		}
	}
}

// Attr ...
func (sr *SoftResource) Attr(key string) Attr {
	sr.check()
	for i := range sr.attrs {
		if sr.attrs[i].Name == key {
			return sr.attrs[i]
		}
	}
	return Attr{}
}

// Rel ...
func (sr *SoftResource) Rel(key string) Rel {
	sr.check()
	for i := range sr.rels {
		if sr.rels[i].Name == key {
			return sr.rels[i]
		}
	}
	return Rel{}
}

// New ...
func (sr *SoftResource) New() Resource {
	sr.check()
	return &SoftResource{
		typ:   sr.typ,
		attrs: append(sr.attrs[:0:0], sr.attrs...),
		rels:  append(sr.rels[:0:0], sr.rels...),
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
	return sr.typ
}

// Get ...
func (sr *SoftResource) Get(key string) interface{} {
	sr.check()
	for k := range sr.attrs {
		if sr.attrs[k].Name == key {
			return sr.data[key]
		}
	}
	return nil
}

// SetID ...
func (sr *SoftResource) SetID(id string) {
	sr.check()
	sr.id = id
}

// SetType ...
func (sr *SoftResource) SetType(typ string) {
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
	for i := range sr.rels {
		if sr.rels[i].ToOne && sr.rels[i].Name == key {
			return sr.data[key].(string)
		}
	}
	return ""
}

// GetToMany ...
func (sr *SoftResource) GetToMany(key string) []string {
	sr.check()
	for i := range sr.rels {
		if !sr.rels[i].ToOne && sr.rels[i].Name == key {
			return sr.data[key].([]string)
		}
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
		typ:   sr.typ,
		id:    sr.id,
		attrs: append(sr.attrs[:0:0], sr.attrs...),
		rels:  append(sr.rels[:0:0], sr.rels...),
		data:  copystructure.Must(copystructure.Copy(sr.data)).(map[string]interface{}),
	}
}

// UnmarshalJSON ...
func (sr *SoftResource) UnmarshalJSON(payload []byte) error {
	sr.check()
	// TODO
	return nil
}

func (sr *SoftResource) fields() []string {
	fields := make([]string, 0, len(sr.attrs)+len(sr.rels))
	for i := range sr.attrs {
		fields = append(fields, sr.attrs[i].Name)
	}
	for i := range sr.rels {
		fields = append(fields, sr.rels[i].Name)
	}
	return fields
}

func (sr *SoftResource) check() {
	if sr.attrs == nil {
		sr.attrs = []Attr{}
	}
	if sr.rels == nil {
		sr.rels = []Rel{}
	}
	if sr.data == nil {
		sr.data = map[string]interface{}{}
	}

	for i := range sr.attrs {
		n := sr.attrs[i].Name
		if _, ok := sr.data[n]; !ok {
			sr.data[n] = ZeroValue(sr.attrs[i].Type, sr.attrs[i].Null)
		}
	}
	for i := range sr.rels {
		n := sr.rels[i].Name
		if _, ok := sr.data[n]; !ok {
			if sr.rels[i].ToOne {
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
