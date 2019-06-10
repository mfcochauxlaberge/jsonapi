package jsonapi

import (
	"fmt"
	"reflect"
	"sort"
)

// Resource ...
type Resource interface {
	// Structure
	Attrs() map[string]Attr
	Rels() map[string]Rel
	Attr(key string) Attr
	Rel(key string) Rel
	New() Resource

	// Read
	GetID() string
	GetType() string
	Get(key string) interface{}

	// Update
	SetID(id string)
	Set(key string, val interface{})

	// Read relationship
	GetToOne(key string) string
	GetToMany(key string) []string

	// Update relationship
	SetToOne(key string, rel string)
	SetToMany(key string, rels []string)

	// Validate
	Validate() []error

	// Copy
	Copy() Resource

	// JSON
	UnmarshalJSON(payload []byte) error
}

// Equal reports whether r1 and r2 are equal.
//
// Two resources are equal if their types are equal, all the attributes
// are equal (same type and same value), and all the relationstips are
// equal.
//
// IDs are ignored.
func Equal(r1, r2 Resource) bool {
	// Type
	if r1.GetType() != r2.GetType() {
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
			// TODO Fix the following condition one day, there should be a better
			// way to do this. Basically, all nils (nil pointer, nil slice, etc)
			// should be considered equal to a nil empty interface.
			if fmt.Sprintf("%v", r1.Get(attr1.Name)) == "<nil>" && fmt.Sprintf("%v", r2.Get(attr1.Name)) == "<nil>" {
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
		return r1Rels[i].Name < r1Rels[j].Name
	})
	rels = r1.Rels()
	r2Rels := make([]Rel, 0, len(rels))
	for name := range rels {
		r2Rels = append(r2Rels, rels[name])
	}
	sort.Slice(r2Rels, func(i, j int) bool {
		return r2Rels[i].Name < r2Rels[j].Name
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
			if r1.GetToOne(rel1.Name) != r2.GetToOne(rel2.Name) {
				return false
			}
		} else {
			if !reflect.DeepEqual(r1.GetToMany(rel1.Name), r2.GetToMany(rel2.Name)) {
				return false
			}
		}
	}

	return true
}

// EqualStrict is like Equal, but it also considers IDs.
func EqualStrict(r1, r2 Resource) bool {
	if r1.GetID() != r2.GetID() {
		return false
	}
	return Equal(r1, r2)
}
