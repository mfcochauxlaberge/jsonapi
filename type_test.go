package jsonapi_test

import (
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	assert := assert.New(t)

	typ := &Type{
		Name: "type1",
	}
	attr1 := Attr{
		Name:     "attr1",
		Type:     AttrTypeString,
		Nullable: false,
	}
	err := typ.AddAttr(attr1)
	assert.NoError(err)
	rel1 := Rel{
		Name: "rel1",
		Type: "type1",
	}
	err = typ.AddRel(rel1)
	assert.NoError(err)

	assert.Contains(typ.Attrs, "attr1")
	assert.Contains(typ.Rels, "rel1")

	// Add invalid attribute (no name)
	err = typ.AddAttr(Attr{})
	assert.Error(err)

	// Add invalid attribute (invalid type)
	err = typ.AddAttr(Attr{Name: "invalid"})
	assert.Error(err)

	// Add invalid attribute (name already used)
	err = typ.AddAttr(Attr{Name: "attr1", Type: AttrTypeString})
	assert.Error(err)

	// Add invalid relationship (no name)
	err = typ.AddRel(Rel{})
	assert.Error(err)

	// Add invalid relationship (empty type)
	err = typ.AddRel(Rel{Name: "invalid"})
	assert.Error(err)

	// Add invalid relationship (name already used)
	err = typ.AddRel(Rel{Name: "rel1", Type: "type1"})
	assert.Error(err)
}

func TestTypeNewFunc(t *testing.T) {
	assert := assert.New(t)

	// NewFunc is nil
	typ := &Type{}
	assert.Equal(&SoftResource{Type: typ}, typ.New())

	// NewFunc is not nil
	typ = &Type{
		NewFunc: func() Resource {
			res := &SoftResource{}
			res.SetID("abc123")
			return res
		},
	}
	assert.Equal("abc123", typ.New().GetID())
}

func TestInverseRel(t *testing.T) {
	assert := assert.New(t)

	rel := Rel{
		Name:         "rel1",
		Type:         "type2",
		ToOne:        true,
		InverseName:  "rel2",
		InverseType:  "type1",
		InverseToOne: false,
	}

	invRel := rel.Inverse()

	assert.Equal("rel2", invRel.Name)
	assert.Equal("type1", invRel.Type)
	assert.Equal(false, invRel.ToOne)
	assert.Equal("rel1", invRel.InverseName)
	assert.Equal("type2", invRel.InverseType)
	assert.Equal(true, invRel.InverseToOne)
}

func TestGetAttrType(t *testing.T) {
	assert := assert.New(t)

	typ, nullable := GetAttrType("string")
	assert.Equal(AttrTypeString, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("int")
	assert.Equal(AttrTypeInt, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("int8")
	assert.Equal(AttrTypeInt8, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("int16")
	assert.Equal(AttrTypeInt16, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("int32")
	assert.Equal(AttrTypeInt32, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("int64")
	assert.Equal(AttrTypeInt64, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("uint")
	assert.Equal(AttrTypeUint, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("uint8")
	assert.Equal(AttrTypeUint8, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("uint16")
	assert.Equal(AttrTypeUint16, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("uint32")
	assert.Equal(AttrTypeUint32, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("uint64")
	assert.Equal(AttrTypeUint64, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("bool")
	assert.Equal(AttrTypeBool, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("time.Time")
	assert.Equal(AttrTypeTime, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("*string")
	assert.Equal(AttrTypeString, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*int")
	assert.Equal(AttrTypeInt, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*int8")
	assert.Equal(AttrTypeInt8, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*int16")
	assert.Equal(AttrTypeInt16, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*int32")
	assert.Equal(AttrTypeInt32, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*int64")
	assert.Equal(AttrTypeInt64, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*uint")
	assert.Equal(AttrTypeUint, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*uint8")
	assert.Equal(AttrTypeUint8, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*uint16")
	assert.Equal(AttrTypeUint16, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*uint32")
	assert.Equal(AttrTypeUint32, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*uint64")
	assert.Equal(AttrTypeUint64, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*bool")
	assert.Equal(AttrTypeBool, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*time.Time")
	assert.Equal(AttrTypeTime, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("invalid")
	assert.Equal(AttrTypeInvalid, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("")
	assert.Equal(AttrTypeInvalid, typ)
	assert.False(nullable)
}

func TestGetAttrTypeString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("string", GetAttrTypeString(AttrTypeString, false))
	assert.Equal("int", GetAttrTypeString(AttrTypeInt, false))
	assert.Equal("int8", GetAttrTypeString(AttrTypeInt8, false))
	assert.Equal("int16", GetAttrTypeString(AttrTypeInt16, false))
	assert.Equal("int32", GetAttrTypeString(AttrTypeInt32, false))
	assert.Equal("int64", GetAttrTypeString(AttrTypeInt64, false))
	assert.Equal("uint", GetAttrTypeString(AttrTypeUint, false))
	assert.Equal("uint8", GetAttrTypeString(AttrTypeUint8, false))
	assert.Equal("uint16", GetAttrTypeString(AttrTypeUint16, false))
	assert.Equal("uint32", GetAttrTypeString(AttrTypeUint32, false))
	assert.Equal("uint64", GetAttrTypeString(AttrTypeUint64, false))
	assert.Equal("bool", GetAttrTypeString(AttrTypeBool, false))
	assert.Equal("time.Time", GetAttrTypeString(AttrTypeTime, false))
	assert.Equal("*string", GetAttrTypeString(AttrTypeString, true))
	assert.Equal("*int", GetAttrTypeString(AttrTypeInt, true))
	assert.Equal("*int8", GetAttrTypeString(AttrTypeInt8, true))
	assert.Equal("*int16", GetAttrTypeString(AttrTypeInt16, true))
	assert.Equal("*int32", GetAttrTypeString(AttrTypeInt32, true))
	assert.Equal("*int64", GetAttrTypeString(AttrTypeInt64, true))
	assert.Equal("*uint", GetAttrTypeString(AttrTypeUint, true))
	assert.Equal("*uint8", GetAttrTypeString(AttrTypeUint8, true))
	assert.Equal("*uint16", GetAttrTypeString(AttrTypeUint16, true))
	assert.Equal("*uint32", GetAttrTypeString(AttrTypeUint32, true))
	assert.Equal("*uint64", GetAttrTypeString(AttrTypeUint64, true))
	assert.Equal("*bool", GetAttrTypeString(AttrTypeBool, true))
	assert.Equal("*time.Time", GetAttrTypeString(AttrTypeTime, true))
	assert.Equal("", GetAttrTypeString(AttrTypeInvalid, false))
	assert.Equal("", GetAttrTypeString(999, false))
}

func TestGetZeroValue(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("", GetZeroValue(AttrTypeString, false))
	assert.Equal(int(0), GetZeroValue(AttrTypeInt, false))
	assert.Equal(int8(0), GetZeroValue(AttrTypeInt8, false))
	assert.Equal(int16(0), GetZeroValue(AttrTypeInt16, false))
	assert.Equal(int32(0), GetZeroValue(AttrTypeInt32, false))
	assert.Equal(int64(0), GetZeroValue(AttrTypeInt64, false))
	assert.Equal(uint(0), GetZeroValue(AttrTypeUint, false))
	assert.Equal(uint8(0), GetZeroValue(AttrTypeUint8, false))
	assert.Equal(uint16(0), GetZeroValue(AttrTypeUint16, false))
	assert.Equal(uint32(0), GetZeroValue(AttrTypeUint32, false))
	assert.Equal(uint64(0), GetZeroValue(AttrTypeUint64, false))
	assert.Equal(false, GetZeroValue(AttrTypeBool, false))
	assert.Equal(time.Time{}, GetZeroValue(AttrTypeTime, false))
	assert.Equal(nilptr("string"), GetZeroValue(AttrTypeString, true))
	assert.Equal(nilptr("int"), GetZeroValue(AttrTypeInt, true))
	assert.Equal(nilptr("int8"), GetZeroValue(AttrTypeInt8, true))
	assert.Equal(nilptr("int16"), GetZeroValue(AttrTypeInt16, true))
	assert.Equal(nilptr("int32"), GetZeroValue(AttrTypeInt32, true))
	assert.Equal(nilptr("int64"), GetZeroValue(AttrTypeInt64, true))
	assert.Equal(nilptr("uint"), GetZeroValue(AttrTypeUint, true))
	assert.Equal(nilptr("uint8"), GetZeroValue(AttrTypeUint8, true))
	assert.Equal(nilptr("uint16"), GetZeroValue(AttrTypeUint16, true))
	assert.Equal(nilptr("uint32"), GetZeroValue(AttrTypeUint32, true))
	assert.Equal(nilptr("uint64"), GetZeroValue(AttrTypeUint64, true))
	assert.Equal(nilptr("bool"), GetZeroValue(AttrTypeBool, true))
	assert.Equal(nilptr("time.Time"), GetZeroValue(AttrTypeTime, true))
	assert.Equal(nil, GetZeroValue(AttrTypeInvalid, false))
	assert.Equal(nil, GetZeroValue(999, false))
}

func TestCopyType(t *testing.T) {
	assert := assert.New(t)

	typ1 := Type{
		Name: "type1",
		Attrs: map[string]Attr{
			"attr1": {
				Name:     "attr1",
				Type:     AttrTypeString,
				Nullable: true,
			},
		},
		Rels: map[string]Rel{
			"rel1": {
				Name:         "rel1",
				Type:         "type2",
				ToOne:        true,
				InverseName:  "rel2",
				InverseType:  "type1",
				InverseToOne: false,
			},
		},
	}

	// Copy
	typ2 := CopyType(typ1)

	assert.Equal("type1", typ2.Name)
	assert.Len(typ2.Attrs, 1)
	assert.Equal("attr1", typ2.Attrs["attr1"].Name)
	assert.Equal(AttrTypeString, typ2.Attrs["attr1"].Type)
	assert.True(typ2.Attrs["attr1"].Nullable)
	assert.Len(typ2.Rels, 1)
	assert.Equal("rel1", typ2.Rels["rel1"].Name)
	assert.Equal("type2", typ2.Rels["rel1"].Type)
	assert.True(typ2.Rels["rel1"].ToOne)
	assert.Equal("rel2", typ2.Rels["rel1"].InverseName)
	assert.Equal("type1", typ2.Rels["rel1"].InverseType)
	assert.False(typ2.Rels["rel1"].InverseToOne)

	// Modify original (copy should not change)
	typ1.Name = "type3"
	typ1.Attrs["attr2"] = Attr{
		Name: "attr2",
		Type: AttrTypeInt,
	}

	assert.Equal("type1", typ2.Name)
	assert.Len(typ2.Attrs, 1)

	typ1.Name = "type1"
	delete(typ1.Attrs, "attr2")

	// Modify copy (original should not change)
	typ2.Name = "type3"
	typ2.Attrs["attr2"] = Attr{
		Name: "attr2",
		Type: AttrTypeInt,
	}

	assert.Equal("type1", typ1.Name)
	assert.Len(typ1.Attrs, 1)
}
