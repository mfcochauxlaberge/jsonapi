package jsonapi_test

import (
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalPartialResource(t *testing.T) {
	// Setup
	typ, _ := BuildType(mocktype{})
	typ.NewFunc = func() Resource {
		return Wrap(&mocktype{})
	}
	schema := &Schema{Types: []Type{typ}}

	// Tests
	t.Run("partial resource", func(t *testing.T) {
		assert := assert.New(t)

		payload := `{
			"id": "abc123",
			"type": "mocktype",
			"attributes": {
				"str": "abc"
			},
			"relationships": {
				"to-1": {
					"data": {
						"type": "mocktype",
						"data": "def"
					}
				},
				"to-x": {
					"data": [
						{
							"type": "mocktype",
							"data": "ghi"
						},
						{
							"type": "mocktype",
							"data": "jkl"
						}
					]
				}
			}
		}`

		res, err := UnmarshalPartialResource([]byte(payload), schema)
		assert.NoError(err)

		assert.Equal("abc123", res.GetID())
		assert.Equal("mocktype", res.GetType().Name)
		assert.Len(res.Attrs(), 1)
		assert.Len(res.Rels(), 2)
	})

	t.Run("partial resource (invalid attribute)", func(t *testing.T) {
		assert := assert.New(t)

		payload := `{
			"id": "abc123",
			"type": "mocktype",
			"attributes": {
				"int": "not an int"
			}
		}`

		_, err := UnmarshalPartialResource([]byte(payload), schema)
		assert.EqualError(
			err,
			"400 Bad Request: The field value is invalid for the expected type.",
		)
	})

	t.Run("partial resource (unknown attribute)", func(t *testing.T) {
		assert := assert.New(t)

		payload := `{
			"id": "abc123",
			"type": "mocktype",
			"attributes": {
				"unknown": "abc"
			}
		}`

		_, err := UnmarshalPartialResource([]byte(payload), schema)
		assert.EqualError(
			err,
			`400 Bad Request: "unknown" is not a known field.`,
		)
	})

	t.Run("partial resource (invalid relationship)", func(t *testing.T) {
		assert := assert.New(t)

		payload := `{
			"id": "abc123",
			"type": "mocktype",
			"relationships": {
				"to-1": {
					"data": [
						{
							"type": "mocktype",
							"data": "def"
						}
					]
				}
			}
		}`

		_, err := UnmarshalPartialResource([]byte(payload), schema)
		assert.EqualError(
			err,
			"400 Bad Request: The field value is invalid for the expected type.",
		)
	})

	t.Run("partial resource (unknown relationship)", func(t *testing.T) {
		assert := assert.New(t)

		payload := `{
			"id": "abc123",
			"type": "mocktype",
			"relationships": {
				"unknown": {
					"data": {
						"type": "mocktype",
						"data": "def"
					}
				}
			}
		}`

		_, err := UnmarshalPartialResource([]byte(payload), schema)
		assert.EqualError(
			err,
			`400 Bad Request: "unknown" is not a known field.`,
		)
	})

	t.Run("partial resource (invalid)", func(t *testing.T) {
		assert := assert.New(t)

		payload := `{invalid}`

		_, err := UnmarshalPartialResource([]byte(payload), schema)
		assert.EqualError(
			err,
			"400 Bad Request: The provided JSON body could not be read.",
		)
	})
}

func TestEqual(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()

	mt11 := Wrap(&mockType1{
		ID:             "mt1",
		Str:            "str",
		Int:            1,
		Int8:           2,
		Int16:          3,
		Int32:          4,
		Int64:          5,
		Uint:           6,
		Uint8:          7,
		Uint16:         8,
		Uint32:         9,
		Bool:           true,
		Time:           now,
		ToOne:          "a",
		ToOneFromOne:   "b",
		ToOneFromMany:  "c",
		ToMany:         []string{"a", "b", "c"},
		ToManyFromOne:  []string{"a", "b", "c"},
		ToManyFromMany: []string{"a", "b", "c"},
	})

	mt12 := Wrap(&mockType1{
		ID:             "mt2",
		Str:            "str",
		Int:            1,
		Int8:           2,
		Int16:          3,
		Int32:          4,
		Int64:          5,
		Uint:           6,
		Uint8:          7,
		Uint16:         8,
		Uint32:         9,
		Bool:           true,
		Time:           now,
		ToOne:          "a",
		ToOneFromOne:   "b",
		ToOneFromMany:  "c",
		ToMany:         []string{"a", "b", "c"},
		ToManyFromOne:  []string{"a", "b", "c"},
		ToManyFromMany: []string{"a", "b", "c"},
	})

	mt13 := Wrap(&mockType1{
		ID:             "mt3",
		Str:            "str",
		Int:            11,
		Int8:           12,
		Int16:          13,
		Int32:          14,
		Int64:          15,
		Uint:           16,
		Uint8:          17,
		Uint16:         18,
		Uint32:         19,
		Bool:           false,
		Time:           time.Now(),
		ToOne:          "d",
		ToOneFromOne:   "e",
		ToOneFromMany:  "f",
		ToMany:         []string{"d", "e", "f"},
		ToManyFromOne:  []string{"d", "e", "f"},
		ToManyFromMany: []string{"d", "e", "f"},
	})

	mt21 := Wrap(&mockType2{
		ID:             "mt1",
		StrPtr:         func() *string { v := "id"; return &v }(),
		IntPtr:         func() *int { v := int(1); return &v }(),
		Int8Ptr:        func() *int8 { v := int8(2); return &v }(),
		Int16Ptr:       func() *int16 { v := int16(3); return &v }(),
		Int32Ptr:       func() *int32 { v := int32(4); return &v }(),
		Int64Ptr:       func() *int64 { v := int64(5); return &v }(),
		UintPtr:        func() *uint { v := uint(6); return &v }(),
		Uint8Ptr:       func() *uint8 { v := uint8(7); return &v }(),
		Uint16Ptr:      func() *uint16 { v := uint16(8); return &v }(),
		Uint32Ptr:      func() *uint32 { v := uint32(9); return &v }(),
		BoolPtr:        func() *bool { v := true; return &v }(),
		TimePtr:        func() *time.Time { v := time.Now(); return &v }(),
		ToOneFromOne:   "a",
		ToOneFromMany:  "b",
		ToManyFromOne:  []string{"a", "b", "c"},
		ToManyFromMany: []string{"a", "b", "c"},
	})

	assert.True(Equal(mt11, mt11), "same instance")
	assert.True(Equal(mt11, mt12), "identical resources")
	assert.False(EqualStrict(mt11, mt12), "different IDs")
	assert.False(Equal(mt11, mt13), "different resources (same type)")
	assert.False(Equal(mt11, mt21), "different types")

	typ := mt11.GetType().Copy()
	sr1 := &SoftResource{Type: &typ}
	sr1.RemoveField("str")
	assert.False(Equal(mt11, sr1), "different number of attributes")

	typ = mt11.GetType().Copy()
	sr1 = &SoftResource{Type: &typ}

	for _, attr := range typ.Attrs {
		sr1.Set(attr.Name, mt11.Get(attr.Name))
	}

	for _, rel := range typ.Rels {
		if rel.ToOne {
			sr1.Set(rel.FromName, mt11.Get(rel.FromName).(string))
		} else {
			sr1.Set(rel.FromName, mt11.Get(rel.FromName).([]string))
		}
	}

	sr1.RemoveField("to-one")
	assert.False(Equal(mt11, sr1), "different number of relationships")

	sr1.AddRel(Rel{
		FromName: "to-one",
		ToOne:    false,
		ToType:   "mocktypes2",
	})
	assert.False(Equal(mt11, sr1), "different to-one property")

	sr1.RemoveField("to-one")
	sr1.AddRel(Rel{
		FromName: "to-one",
		ToOne:    true,
		ToType:   "mocktypes2",
	})
	sr1.Set("to-one", "b")
	assert.False(Equal(mt11, sr1), "different relationship value (to-one)")

	sr1.Set("to-one", "a")
	sr1.Set("to-many", []string{"d", "e", "f"})
	assert.False(Equal(mt11, sr1), "different relationship value (to-many)")

	// Comparing two nil values of different types
	sr3 := &SoftResource{}
	sr3.AddAttr(Attr{
		Name:     "nil",
		Type:     AttrTypeString,
		Nullable: true,
	})
	sr3.Set("nil", (*string)(nil))

	sr4 := &SoftResource{}
	sr4.AddAttr(Attr{
		Name:     "nil2",
		Type:     AttrTypeInt,
		Nullable: true,
	})
	sr3.Set("nil", (*int)(nil))

	assert.Equal(true, Equal(sr3, sr4))
}

func TestEqualStrict(t *testing.T) {
	assert := assert.New(t)

	sr1 := &SoftResource{}
	sr1.SetType(&Type{
		Name: "type",
	})

	sr2 := &SoftResource{}
	sr2.SetType(&Type{
		Name: "type",
	})

	// Same ID
	sr1.SetID("an-id")
	sr2.SetID("an-id")
	assert.True(Equal(sr1, sr2))
	assert.True(EqualStrict(sr1, sr2))

	// Different ID
	sr2.SetID("another-id")
	assert.True(Equal(sr1, sr2))
	assert.False(EqualStrict(sr1, sr2))
}
