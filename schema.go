package jsonapi

import (
	"errors"
	"fmt"
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

// Schema ...
type Schema struct {
	Types []Type
}

// AddType ...
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

// RemoveType ...
func (s *Schema) RemoveType(typ string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			s.Types = append(s.Types[0:i], s.Types[i+1:]...)
		}
	}

	return nil
}

// AddAttr ...
func (s *Schema) AddAttr(typ string, attr Attr) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].AddAttr(attr)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// RemoveAttr ...
func (s *Schema) RemoveAttr(typ string, attr string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].RemoveAttr(attr)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// AddRel ...
func (s *Schema) AddRel(typ string, rel Rel) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].AddRel(rel)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// RemoveRel ...
func (s *Schema) RemoveRel(typ string, rel string) error {
	for i := range s.Types {
		if s.Types[i].Name == typ {
			return s.Types[i].RemoveRel(rel)
		}
	}

	return fmt.Errorf("jsonapi: type %s does not exist", typ)
}

// HasType ...
func (s *Schema) HasType(name string) bool {
	for i := range s.Types {
		if s.Types[i].Name == name {
			return true
		}
	}
	return false
}

// GetType ...
func (s *Schema) GetType(name string) (Type, bool) {
	for _, typ := range s.Types {
		if typ.Name == name {
			return typ, true
		}
	}
	return Type{}, false
}

// GetResource ...
func (s *Schema) GetResource(name string) Resource {
	typ, ok := s.GetType(name)
	if ok {
		return NewSoftResource(typ, nil)
	}
	return nil
}

// Check ...
func (s *Schema) Check() []error {
	var (
		ok   bool
		errs = []error{}
	)

	// Check the inverse relationships
	for _, typ := range s.Types {
		// Relationships
		for _, rel := range typ.Rels {
			var targetType Type

			// Does the relationship point to a type that exists?
			if targetType, ok = s.GetType(rel.Type); !ok {
				errs = append(errs, fmt.Errorf(
					"jsonapi: the target type of relationship %s of type %s does not exist",
					rel.Name,
					typ.Name,
				))
			}

			// Inverse relationship (if relevant)
			if rel.InverseName != "" {
				// Is the inverse relationship type the same as its type name?
				if rel.InverseType != typ.Name {
					errs = append(errs, fmt.Errorf(
						"jsonapi: the inverse type of relationship %s should its type's name (%s, not %s)",
						rel.Name,
						typ.Name,
						rel.InverseType,
					))
				}

				// Do both relationships (current and inverse) point to each other?
				var found bool
				for _, invRel := range targetType.Rels {
					if rel.Name == invRel.InverseName && rel.InverseName == invRel.Name {
						found = true
					}
				}
				if !found {
					errs = append(errs, fmt.Errorf(
						"jsonapi: relationship %s of type %s and its inverse do not point each other",
						rel.Name,
						typ.Name,
					))
				}
			}

		}
	}

	return errs
}

// GetAttrType ...
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

// GetAttrTypeString ...
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

// GetZeroValue ...
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
