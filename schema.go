package jsonapi

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

// Attribute types
const (
	AttrTypeInvalid = iota
	AttrTypeString
	AttrTypeInt
	AttrTypeInt8
	AttrTypeInt16
	AttrTypeInt32
	AttrTypeInt64
	AttrTypeUint
	AttrTypeUint8
	AttrTypeUint16
	AttrTypeUint32
	AttrTypeBool
	AttrTypeTime
)

// A Schema contains a list of types. It makes sure that each type is valid and
// unique.
//
// Check can be used to validate the relationships between the types.
type Schema struct {
	Types []Type
}

// AddType adds a type to the schema.
func (s *Schema) AddType(typ Type) error {
	// Validation
	if typ.Name == "" {
		return errors.New("jsonapi: type name is empty")
	}

	// Make sure the name isn't already used
	for i := range s.Types {
		if s.Types[i].Name == typ.Name {
			return fmt.Errorf("jsonapi: type name %s is already used", typ.Name)
		}
	}

	s.Types = append(s.Types, typ)

	return nil
}

// RemoveType removes a type from the schema.
func (s *Schema) RemoveType(typ string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			s.Types = append(s.Types[0:i], s.Types[i+1:]...)
		}
	}

	return nil
}

// AddAttr adds an attribute to the specified type.
func (s *Schema) AddAttr(typ string, attr Attr) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].AddAttr(attr)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// RemoveAttr removes an attribute from the specified type.
func (s *Schema) RemoveAttr(typ string, attr string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].RemoveAttr(attr)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// AddRel adds a relationship to the specified type.
func (s *Schema) AddRel(typ string, rel Rel) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].AddRel(rel)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// RemoveRel removes a relationship from the specified type.
func (s *Schema) RemoveRel(typ string, rel string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].RemoveRel(rel)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// HasType returns a boolean indicating whether a type has the specified name
// or not.
func (s *Schema) HasType(name string) bool {
	for i := range s.Types {
		if s.Types[i].Name == name {
			return true
		}
	}
	return false
}

// GetType returns the type associated with the speficied name.
func (s *Schema) GetType(name string) (Type, bool) {
	var typ Type
	for _, typ = range s.Types {
		if typ.Name == name {
			break
		}
	}
	return typ, false
}

// Check checks the integrity of all the relationships between the types and
// returns all the errors that were found.
func (s *Schema) Check() []error {
	// TODO Don't use Registry (which should be removed anyway)
	reg := Registry{}
	for _, typ := range s.Types {
		reg.Types[typ.Name] = typ
	}
	return reg.Check()
}

// Type represents a JSON:API type.
type Type struct {
	Name    string
	Attrs   map[string]Attr
	Rels    map[string]Rel
	Default Resource
}

// AddAttr adds an attribute to the type.
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

// RemoveAttr removes an attribute from the type.
func (t *Type) RemoveAttr(attr string) error {
	for i := range t.Attrs {
		if t.Attrs[i].Name == attr {
			delete(t.Attrs, attr)
		}
	}

	return nil
}

// AddRel adds a relationship to the type.
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

// RemoveRel removes a relationship from the type.
func (t *Type) RemoveRel(rel string) error {
	for i := range t.Rels {
		if t.Rels[i].Name == rel {
			delete(t.Rels, rel)
		}
	}

	return nil
}

// Fields returns a slice of strings that contains the name of each attribute
// and relationship combined. The list is sorted alphabetically.
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

// Attr represents a JSON:API attribute.
type Attr struct {
	Name string
	Type int
	Null bool
}

// Rel represent a JSON:API relationship between two types.
//
// Name is the name of the field. Type is the name of the type this
// relationship points to. ToOne is true for to-one relationships and
// false for to-many relationships.
//
// For two-way relationships, the inverse fields must be set to their
// corresponding values from the opposite relationship. For one-way
// relationships, the values must be zeroed.
type Rel struct {
	Name         string
	Type         string
	ToOne        bool
	InverseName  string
	InverseType  string
	InverseToOne bool
}

// Inverse returns the opposite relationship.
//
// If the relationship is only one way, Name, Type, and ToOne will be set to
// their zero values.
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

// GetAttrType returns the attribute type as an int (see constants) and a
// boolean that indicates whether the attribute can be null or not.
func GetAttrType(t string) (int, bool) {
	t2 := t
	if strings.HasPrefix(t2, "*") {
		t2 = t[1:]
	}
	switch t2 {
	case "string":
		if t[0] != '*' {
			return AttrTypeString, false
		}
		return AttrTypeString, true
	case "int":
		if t[0] != '*' {
			return AttrTypeInt, false
		}
		return AttrTypeInt, true
	case "int8":
		if t[0] != '*' {
			return AttrTypeInt8, false
		}
		return AttrTypeInt8, true
	case "int16":
		if t[0] != '*' {
			return AttrTypeInt16, false
		}
		return AttrTypeInt16, true
	case "int32":
		if t[0] != '*' {
			return AttrTypeInt32, false
		}
		return AttrTypeInt32, true
	case "int64":
		if t[0] != '*' {
			return AttrTypeInt64, false
		}
		return AttrTypeInt64, true
	case "uint":
		if t[0] != '*' {
			return AttrTypeUint, false
		}
		return AttrTypeUint, true
	case "uint8":
		if t[0] != '*' {
			return AttrTypeUint8, false
		}
		return AttrTypeUint8, true
	case "uint16":
		if t[0] != '*' {
			return AttrTypeUint16, false
		}
		return AttrTypeUint16, true
	case "uint32":
		if t[0] != '*' {
			return AttrTypeUint32, false
		}
		return AttrTypeUint32, true
	case "bool":
		if t[0] != '*' {
			return AttrTypeBool, false
		}
		return AttrTypeBool, true
	case "time.Time":
		if t[0] != '*' {
			return AttrTypeTime, false
		}
		return AttrTypeTime, true
	default:
		if t[0] != '*' {
			return AttrTypeInvalid, false
		}
		return AttrTypeInvalid, true
	}
}

// GetAttrTypeString return the name of the attribute type specified by an int
// (see constants) and a boolean that indicates whether the value can be null
// or not.
func GetAttrTypeString(t int, null bool) string {
	switch t {
	case AttrTypeString:
		if !null {
			return "string"
		}
		return "*string"
	case AttrTypeInt:
		if !null {
			return "int"
		}
		return "*int"
	case AttrTypeInt8:
		if !null {
			return "int8"
		}
		return "*int8"
	case AttrTypeInt16:
		if !null {
			return "int16"
		}
		return "*int16"
	case AttrTypeInt32:
		if !null {
			return "int32"
		}
		return "*int32"
	case AttrTypeInt64:
		if !null {
			return "int64"
		}
		return "*int64"
	case AttrTypeUint:
		if !null {
			return "uint"
		}
		return "*uint"
	case AttrTypeUint8:
		if !null {
			return "uint8"
		}
		return "*uint8"
	case AttrTypeUint16:
		if !null {
			return "uint16"
		}
		return "*uint16"
	case AttrTypeUint32:
		if !null {
			return "uint32"
		}
		return "*uint32"
	case AttrTypeBool:
		if !null {
			return "bool"
		}
		return "*bool"
	case AttrTypeTime:
		if !null {
			return "time"
		}
		return "*time.Time"
	default:
		return ""
	}
}

// GetZeroValue returns the zero value of the attribute type represented by the
// specified int (see constants).
//
// If null is true, the returned value is a nil pointer.
func GetZeroValue(t int, null bool) interface{} {
	switch t {
	case AttrTypeString:
		v := ""
		if !null {
			return v
		}
		return &v
	case AttrTypeInt:
		v := int(0)
		if !null {
			return v
		}
		return &v
	case AttrTypeInt8:
		v := int8(0)
		if !null {
			return v
		}
		return &v
	case AttrTypeInt16:
		v := int16(0)
		if !null {
			return v
		}
		return &v
	case AttrTypeInt32:
		v := int32(0)
		if !null {
			return v
		}
		return &v
	case AttrTypeInt64:
		v := int64(0)
		if !null {
			return v
		}
		return &v
	case AttrTypeUint:
		v := uint(0)
		if !null {
			return v
		}
		return &v
	case AttrTypeUint8:
		v := uint8(0)
		if !null {
			return v
		}
		return &v
	case AttrTypeUint16:
		v := uint16(0)
		if !null {
			return v
		}
		return &v
	case AttrTypeUint32:
		v := uint32(0)
		if !null {
			return v
		}
		return &v
	case AttrTypeBool:
		v := false
		if !null {
			return v
		}
		return &v
	case AttrTypeTime:
		v := time.Time{}
		if !null {
			return v
		}
		return &v
	default:
		return nil
	}
}
