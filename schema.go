package jsonapi

import (
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

// Type ...
type Type struct {
	Name    string
	Fields  []string
	Attrs   map[string]Attr
	Rels    map[string]Rel
	Default Resource
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

// GetAttrString ...
func GetAttrString(t int, null bool) string {
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

// ZeroValue ...
func ZeroValue(t int, null bool) interface{} {
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
