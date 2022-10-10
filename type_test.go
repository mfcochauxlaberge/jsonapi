package jsonapi_test

import (
	"encoding/json"
	"fmt"
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
		FromName: "rel1",
		ToType:   "type1",
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
	err = typ.AddRel(Rel{FromName: "invalid"})
	assert.Error(err)

	// Add invalid relationship (name already used)
	err = typ.AddRel(Rel{FromName: "rel1", ToType: "type1"})
	assert.Error(err)
}

// TODO Add tests with attributes and relationships.
func TestTypeEqual(t *testing.T) {
	assert := assert.New(t)

	// Two empty types
	typ1 := Type{}
	typ2 := Type{}
	assert.True(typ1.Equal(typ2))

	typ1.Name = "type1"
	typ2.Name = "type1"
	assert.True(typ1.Equal(typ2))

	typ1.Name = "type1"
	typ2.Name = "type2"
	assert.False(typ1.Equal(typ2))

	// Make sure NewFunc is ignored.
	typ1.Name = "type1"
	typ1.NewFunc = func() Resource {
		return nil
	}
	typ2.Name = "type1"
	typ2.NewFunc = func() Resource {
		return &SoftResource{}
	}
	assert.True(typ1.Equal(typ2))
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
	assert.Equal("abc123", typ.New().Get("id").(string))
}

func TestAttrUnmarshalToType(t *testing.T) {
	assert := assert.New(t)

	var (
		vstr    = "str"
		vint    = int(1)
		vint8   = int8(8)
		vint16  = int16(16)
		vint32  = int32(32)
		vint64  = int64(64)
		vuint   = uint(1)
		vuint8  = uint8(8)
		vuint16 = uint16(16)
		vuint32 = uint32(32)
		vuint64 = uint64(64)
		vbool   = true
	)

	tests := []struct {
		val any
	}{
		{val: "str"},            // string
		{val: 1},                // int
		{val: int8(8)},          // int8
		{val: int16(16)},        // int16
		{val: int32(32)},        // int32
		{val: int64(64)},        // int64
		{val: uint(1)},          // uint
		{val: uint8(8)},         // uint8
		{val: uint16(16)},       // uint16
		{val: uint32(32)},       // uint32
		{val: uint64(64)},       // uint64
		{val: true},             // bool
		{val: time.Time{}},      // time
		{val: []byte{1, 2, 3}},  // []byte
		{val: &vstr},            // *string
		{val: &vint},            // *int
		{val: &vint8},           // *int8
		{val: &vint16},          // *int16
		{val: &vint32},          // *int32
		{val: &vint64},          // *int64
		{val: &vuint},           // *uint
		{val: &vuint8},          // *uint8
		{val: &vuint16},         // *uint16
		{val: &vuint32},         // *uint32
		{val: &vuint64},         // *uint64
		{val: &vbool},           // *bool
		{val: &time.Time{}},     // *time
		{val: &[]byte{1, 2, 3}}, // *[]byte
	}

	attr := Attr{}

	for _, test := range tests {
		attr.Type, attr.Nullable = GetAttrType(fmt.Sprintf("%T", test.val))
		p, _ := json.Marshal(test.val)
		val, err := attr.UnmarshalToType(p)
		assert.NoError(err)
		assert.Equal(test.val, val)
		assert.Equal(fmt.Sprintf("%T", test.val), fmt.Sprintf("%T", val))
	}

	// Null value
	attr.Nullable = true
	val, err := attr.UnmarshalToType([]byte("null"))
	assert.NoError(err)
	assert.Nil(val)

	// False value
	attr.Type = AttrTypeBool
	val, err = attr.UnmarshalToType([]byte("nottrue"))
	assert.Error(err)
	assert.Nil(val)

	// Invalid slide of bytes
	attr.Type = AttrTypeBytes

	assert.Panics(func() {
		_, _ = attr.UnmarshalToType([]byte("invalid"))
	})
	// assert.Error(err)
	// assert.Nil(val)

	// Invalid attribute type
	attr.Type = AttrTypeInvalid
	val, err = attr.UnmarshalToType([]byte("invalid"))
	err2, ok := err.(Error)
	assert.True(ok)
	assert.IsType(Error{}, err2)
	assert.Nil(val)
}

func TestRelInvert(t *testing.T) {
	assert := assert.New(t)

	rel := Rel{
		FromName: "rel1",
		FromType: "type1",
		ToOne:    true,
		ToName:   "rel2",
		ToType:   "type2",
		FromOne:  false,
	}

	invRel := rel.Invert()

	assert.Equal("rel2", invRel.FromName)
	assert.Equal("type1", invRel.ToType)
	assert.Equal(false, invRel.ToOne)
	assert.Equal("rel1", invRel.ToName)
	assert.Equal("type2", invRel.FromType)
	assert.Equal(true, invRel.FromOne)
}

func TestRelNormalize(t *testing.T) {
	assert := assert.New(t)

	rel := Rel{
		FromName: "rel2",
		FromType: "type2",
		ToOne:    false,
		ToName:   "rel1",
		ToType:   "type1",
		FromOne:  true,
	}

	// Normalize should return the inverse because
	// type1 comes before type2 alphabetically.
	norm := rel.Normalize()
	assert.Equal("type1", norm.FromType)
	assert.Equal("rel1", norm.FromName)
	assert.Equal(true, norm.ToOne)
	assert.Equal("type2", norm.ToType)
	assert.Equal("rel2", norm.ToName)
	assert.Equal(false, norm.FromOne)

	// Normalize again, but it should stay the same.
	norm = norm.Normalize()
	assert.Equal("type1", norm.FromType)
	assert.Equal("rel1", norm.FromName)
	assert.Equal(true, norm.ToOne)
	assert.Equal("type2", norm.ToType)
	assert.Equal("rel2", norm.ToName)
	assert.Equal(false, norm.FromOne)
}

func TestRelString(t *testing.T) {
	assert := assert.New(t)

	rel := Rel{
		FromName: "rel2",
		FromType: "type2",
		ToOne:    false,
		ToName:   "rel1",
		ToType:   "type1",
		FromOne:  true,
	}

	assert.Equal("type1_rel1_type2_rel2", rel.String())
	assert.Equal("type1_rel1_type2_rel2", rel.Invert().String())
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

	typ, nullable = GetAttrType("time")
	assert.Equal(AttrTypeTime, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("[]uint8")
	assert.Equal(AttrTypeBytes, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("[]byte")
	assert.Equal(AttrTypeBytes, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("bytes")
	assert.Equal(AttrTypeBytes, typ)
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

	typ, nullable = GetAttrType("*time")
	assert.Equal(AttrTypeTime, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*[]uint8")
	assert.Equal(AttrTypeBytes, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*[]byte")
	assert.Equal(AttrTypeBytes, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("*bytes")
	assert.Equal(AttrTypeBytes, typ)
	assert.True(nullable)

	typ, nullable = GetAttrType("invalid")
	assert.Equal(AttrTypeInvalid, typ)
	assert.False(nullable)

	typ, nullable = GetAttrType("")
	assert.Equal(AttrTypeInvalid, typ)
	assert.False(nullable)
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
	assert.Equal([]byte{}, GetZeroValue(AttrTypeBytes, false))
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
	assert.Equal(nilptr("[]byte"), GetZeroValue(AttrTypeBytes, true))
	assert.Equal(nil, GetZeroValue(AttrTypeInvalid, false))
	assert.Equal(nil, GetZeroValue("", false))
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
				FromName: "rel1",
				FromType: "type1",
				ToOne:    true,
				ToName:   "rel2",
				ToType:   "type2",
				FromOne:  false,
			},
		},
	}

	// Copy
	typ2 := typ1.Copy()

	assert.Equal("type1", typ2.Name)
	assert.Len(typ2.Attrs, 1)
	assert.Equal("attr1", typ2.Attrs["attr1"].Name)
	assert.Equal(AttrTypeString, typ2.Attrs["attr1"].Type)
	assert.True(typ2.Attrs["attr1"].Nullable)
	assert.Len(typ2.Rels, 1)
	assert.Equal("rel1", typ2.Rels["rel1"].FromName)
	assert.Equal("type2", typ2.Rels["rel1"].ToType)
	assert.True(typ2.Rels["rel1"].ToOne)
	assert.Equal("rel2", typ2.Rels["rel1"].ToName)
	assert.Equal("type1", typ2.Rels["rel1"].FromType)
	assert.False(typ2.Rels["rel1"].FromOne)

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
