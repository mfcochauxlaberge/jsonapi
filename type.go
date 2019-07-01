package jsonapi

import (
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
	AttrTypeUint64
	AttrTypeBool
	AttrTypeTime
)

// A Type stores all the necessary information about a type as represented
// in the JSON:API specification.
type Type struct {
	Name  string
	Attrs map[string]Attr
	Rels  map[string]Rel
}

// AddAttr adds an attributes to the type.
func (t *Type) AddAttr(attr Attr) error {
	// Validation
	if attr.Name == "" {
		return fmt.Errorf("jsonapi: attribute name is empty")
	}

	if GetAttrTypeString(attr.Type, attr.Nullable) == "" {
		return fmt.Errorf("jsonapi: attribute type is invalid")
	}

	// Make sure the name isn't already used
	for i := range t.Attrs {
		if t.Attrs[i].Name == attr.Name {
			return fmt.Errorf("jsonapi: attribute name %s is already used", attr.Name)
		}
	}

	if t.Attrs == nil {
		t.Attrs = map[string]Attr{}
	}
	t.Attrs[attr.Name] = attr

	return nil
}

// RemoveAttr removes an attribute from the type if it exists.
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

	if t.Rels == nil {
		t.Rels = map[string]Rel{}
	}
	t.Rels[rel.Name] = rel

	return nil
}

// RemoveRel removes a relationship from the type if it exists.
func (t *Type) RemoveRel(rel string) error {
	for i := range t.Rels {
		if t.Rels[i].Name == rel {
			delete(t.Rels, rel)
		}
	}

	return nil
}

// Fields returns a list of the names of all the fields (attributes and
// relationships) in the type.
func (t *Type) Fields() []string {
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

// Attr represents a resource attribute.
type Attr struct {
	Name     string
	Type     int
	Nullable bool
}

// Rel represents a resource relationship.
type Rel struct {
	Name         string
	Type         string
	ToOne        bool
	InverseName  string
	InverseType  string
	InverseToOne bool
}

// Inverse returns the inverse relationship of r.
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

// GetAttrType returns the attribute type as an int (see constants) and
// a boolean that indicates whether the attribute can be null or not.
func GetAttrType(t string) (int, bool) {
	nullable := strings.HasPrefix(t, "*")
	if nullable {
		t = t[1:]
	}
	switch t {
	case "string":
		return AttrTypeString, nullable
	case "int":
		return AttrTypeInt, nullable
	case "int8":
		return AttrTypeInt8, nullable
	case "int16":
		return AttrTypeInt16, nullable
	case "int32":
		return AttrTypeInt32, nullable
	case "int64":
		return AttrTypeInt64, nullable
	case "uint":
		return AttrTypeUint, nullable
	case "uint8":
		return AttrTypeUint8, nullable
	case "uint16":
		return AttrTypeUint16, nullable
	case "uint32":
		return AttrTypeUint32, nullable
	case "uint64":
		return AttrTypeUint64, nullable
	case "bool":
		return AttrTypeBool, nullable
	case "time.Time":
		return AttrTypeTime, nullable
	default:
		return AttrTypeInvalid, false
	}
}

// GetAttrTypeString return the name of the attribute type specified
// by an int (see constants) and a boolean that indicates whether the
// value can be null or not.
func GetAttrTypeString(t int, nullable bool) string {
	str := ""
	switch t {
	case AttrTypeString:
		str = "string"
	case AttrTypeInt:
		str = "int"
	case AttrTypeInt8:
		str = "int8"
	case AttrTypeInt16:
		str = "int16"
	case AttrTypeInt32:
		str = "int32"
	case AttrTypeInt64:
		str = "int64"
	case AttrTypeUint:
		str = "uint"
	case AttrTypeUint8:
		str = "uint8"
	case AttrTypeUint16:
		str = "uint16"
	case AttrTypeUint32:
		str = "uint32"
	case AttrTypeUint64:
		str = "uint64"
	case AttrTypeBool:
		str = "bool"
	case AttrTypeTime:
		str = "time.Time"
	default:
		str = ""
	}
	if nullable {
		return "*" + str
	}
	return str
}

// GetZeroValue returns the zero value of the attribute type represented
// by the specified int (see constants).
//
// If null is true, the returned value is a nil pointer.
func GetZeroValue(t int, null bool) interface{} {
	switch t {
	case AttrTypeString:
		if null {
			var np *string
			return np
		}
		return ""
	case AttrTypeInt:
		if null {
			var np *int
			return np
		}
		return int(0)
	case AttrTypeInt8:
		if null {
			var np *int8
			return np
		}
		return int8(0)
	case AttrTypeInt16:
		if null {
			var np *int16
			return np
		}
		return int16(0)
	case AttrTypeInt32:
		if null {
			var np *int32
			return np
		}
		return int32(0)
	case AttrTypeInt64:
		if null {
			var np *int64
			return np
		}
		return int64(0)
	case AttrTypeUint:
		if null {
			var np *uint
			return np
		}
		return uint(0)
	case AttrTypeUint8:
		if null {
			var np *uint8
			return np
		}
		return uint8(0)
	case AttrTypeUint16:
		if null {
			var np *uint16
			return np
		}
		return uint16(0)
	case AttrTypeUint32:
		if null {
			var np *uint32
			return np
		}
		return uint32(0)
	case AttrTypeUint64:
		if null {
			var np *uint64
			return np
		}
		return uint64(0)
	case AttrTypeBool:
		if null {
			var np *bool
			return np
		}
		return false
	case AttrTypeTime:
		if null {
			var np *time.Time
			return np
		}
		return time.Time{}
	default:
		return nil
	}
}
