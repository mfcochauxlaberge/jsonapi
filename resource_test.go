package jsonapi_test

import (
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"
	"github.com/mfcochauxlaberge/tchek"
)

func TestEqual(t *testing.T) {
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
		StrPtr:         func() *string { v := string(1); return &v }(),
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

	tchek.AreEqual(t, "compare same resource with itself", true, Equal(mt11, mt11))
	tchek.AreEqual(t, "compare two identical resources", true, Equal(mt11, mt12))
	tchek.AreEqual(t, "compare two identical resources (different IDs)", false, EqualStrict(mt11, mt12))
	tchek.AreEqual(t, "compare two different resources", false, Equal(mt11, mt13))
	tchek.AreEqual(t, "compare resources of different types", false, Equal(mt11, mt21))
}
