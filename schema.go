package jsonapi

import "time"

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
	AttrTypeStringPtr
	AttrTypeIntPtr
	AttrTypeInt8Ptr
	AttrTypeInt16Ptr
	AttrTypeInt32Ptr
	AttrTypeInt64Ptr
	AttrTypeUintPtr
	AttrTypeUint8Ptr
	AttrTypeUint16Ptr
	AttrTypeUint32Ptr
	AttrTypeBoolPtr
	AttrTypeTimePtr
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
func GetAttrType(t string) int {
	switch t {
	case "string":
		return AttrTypeString
	case "int":
		return AttrTypeInt
	case "int8":
		return AttrTypeInt8
	case "int16":
		return AttrTypeInt16
	case "int32":
		return AttrTypeInt32
	case "int64":
		return AttrTypeInt64
	case "uint":
		return AttrTypeUint
	case "uint8":
		return AttrTypeUint8
	case "uint16":
		return AttrTypeUint16
	case "uint32":
		return AttrTypeUint32
	case "bool":
		return AttrTypeBool
	case "time.Time":
		return AttrTypeTime
	case "*string":
		return AttrTypeStringPtr
	case "*int":
		return AttrTypeIntPtr
	case "*int8":
		return AttrTypeInt8Ptr
	case "*int16":
		return AttrTypeInt16Ptr
	case "*int32":
		return AttrTypeInt32Ptr
	case "*int64":
		return AttrTypeInt64Ptr
	case "*uint":
		return AttrTypeUintPtr
	case "*uint8":
		return AttrTypeUint8Ptr
	case "*uint16":
		return AttrTypeUint16Ptr
	case "*uint32":
		return AttrTypeUint32Ptr
	case "*bool":
		return AttrTypeBoolPtr
	case "*time.Time":
		return AttrTypeTimePtr
	default:
		return AttrTypeInvalid
	}
}

// GetAttrString ...
func GetAttrString(t int) string {
	switch t {
	case AttrTypeString:
		return "string"
	case AttrTypeInt:
		return "int"
	case AttrTypeInt8:
		return "int8"
	case AttrTypeInt16:
		return "int16"
	case AttrTypeInt32:
		return "int32"
	case AttrTypeInt64:
		return "int64"
	case AttrTypeUint:
		return "uint"
	case AttrTypeUint8:
		return "uint8"
	case AttrTypeUint16:
		return "uint16"
	case AttrTypeUint32:
		return "uint32"
	case AttrTypeBool:
		return "bool"
	case AttrTypeTime:
		return "time.Time"
	case AttrTypeStringPtr:
		return "*string"
	case AttrTypeIntPtr:
		return "*int"
	case AttrTypeInt8Ptr:
		return "*int8"
	case AttrTypeInt16Ptr:
		return "*int16"
	case AttrTypeInt32Ptr:
		return "*int32"
	case AttrTypeInt64Ptr:
		return "*int64"
	case AttrTypeUintPtr:
		return "*uint"
	case AttrTypeUint8Ptr:
		return "*uint8"
	case AttrTypeUint16Ptr:
		return "*uint16"
	case AttrTypeUint32Ptr:
		return "*uint32"
	case AttrTypeBoolPtr:
		return "*bool"
	case AttrTypeTimePtr:
		return "*time.Time"
	default:
		return ""
	}
}

// ZeroValue ...
func ZeroValue(t int) interface{} {
	switch t {
	case AttrTypeString:
		return ""
	case AttrTypeInt:
		return int(0)
	case AttrTypeInt8:
		return int8(0)
	case AttrTypeInt16:
		return int16(0)
	case AttrTypeInt32:
		return int32(0)
	case AttrTypeInt64:
		return int64(0)
	case AttrTypeUint:
		return uint(0)
	case AttrTypeUint8:
		return uint8(0)
	case AttrTypeUint16:
		return uint16(0)
	case AttrTypeUint32:
		return uint32(0)
	case AttrTypeBool:
		return false
	case AttrTypeTime:
		return time.Time{}
	case AttrTypeStringPtr:
		var v *string
		return v
	case AttrTypeIntPtr:
		var v *int
		return v
	case AttrTypeInt8Ptr:
		var v *int8
		return v
	case AttrTypeInt16Ptr:
		var v *int16
		return v
	case AttrTypeInt32Ptr:
		var v *int32
		return v
	case AttrTypeInt64Ptr:
		var v *int64
		return v
	case AttrTypeUintPtr:
		var v *uint
		return v
	case AttrTypeUint8Ptr:
		var v *uint8
		return v
	case AttrTypeUint16Ptr:
		var v *uint16
		return v
	case AttrTypeUint32Ptr:
		var v *uint32
		return v
	case AttrTypeBoolPtr:
		var v *bool
		return v
	case AttrTypeTimePtr:
		var v *time.Time
		return v
	default:
		return ""
	}
}
