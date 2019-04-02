package jsonapi

import (
	"fmt"
	"sort"
)

// Type ...
type Type struct {
	Name  string
	Attrs map[string]Attr
	Rels  map[string]Rel
}

// AddAttr ...
func (t *Type) AddAttr(attr Attr) error {
	// Validation
	if attr.Name == "" {
		return fmt.Errorf("jsonapi: attribute name is empty")
	}

	if GetAttrTypeString(attr.Type, attr.Null) == "" {
		return fmt.Errorf("jsonapi: attribute type is invalid")
	}

	// Make sure the name isn't already used
	for i := range t.Attrs {
		if t.Attrs[i].Name == attr.Name {
			return fmt.Errorf("jsonapi: attribute name %s is already used", attr.Name)
		}
	}

	t.Attrs[attr.Name] = attr

	return nil
}

// RemoveAttr ...
func (t *Type) RemoveAttr(attr string) error {
	for i := range t.Attrs {
		if t.Attrs[i].Name == attr {
			delete(t.Attrs, attr)
		}
	}

	return nil
}

// AddRel ...
func (t *Type) AddRel(rel Rel) error {
	// Validation
	if rel.Name == "" {
		return fmt.Errorf("jsonapi: relationship name is empty")
	}
	if rel.Type == "" {
		return fmt.Errorf("jsonapi: relationship type is empty")
	}

	// Make sure the name isn't already used
	for i := range t.Rels {
		if t.Rels[i].Name == rel.Name {
			return fmt.Errorf("jsonapi: relationship name %s is already used", rel.Name)
		}
	}

	t.Rels[rel.Name] = rel

	return nil
}

// RemoveRel ...
func (t *Type) RemoveRel(rel string) error {
	for i := range t.Rels {
		if t.Rels[i].Name == rel {
			delete(t.Rels, rel)
		}
	}

	return nil
}

// Fields ...
func (t Type) Fields() []string {
	fields := make([]string, 0, len(t.Attrs)+len(t.Rels))
	for i := range t.Attrs {
		fields = append(fields, t.Attrs[i].Name)
	}
	for i := range t.Rels {
		fields = append(fields, t.Rels[i].Name)
	}
	sort.Strings(fields)
	return fields
}

// Attr ...
type Attr struct {
	Name string
	Type int
	Null bool
}

// Rel ...
type Rel struct {
	Name         string
	Type         string
	ToOne        bool
	InverseName  string
	InverseType  string
	InverseToOne bool
}

// Inverse ...
func (r *Rel) Inverse() Rel {
	return Rel{
		Name:         r.InverseName,
		Type:         r.InverseType,
		ToOne:        r.InverseToOne,
		InverseName:  r.Name,
		InverseType:  r.Type,
		InverseToOne: r.ToOne,
	}
}
