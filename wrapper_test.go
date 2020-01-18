package jsonapi_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

var _ Resource = (*Wrapper)(nil)

func TestWrap(t *testing.T) {
	assert := assert.New(t)

	assert.Panics(func() {
		_ = Wrap("just a string")
	}, "panic when not a pointer to a struct")

	assert.Panics(func() {
		str := "just a string"
		_ = Wrap(&str)
	}, "panic when not a pointer to a struct")

	assert.NotPanics(func() {
		_ = Wrap(&mocktype{})
	}, "don't panic when a valid struct")

	assert.NotPanics(func() {
		_ = Wrap(mocktype{})
	}, "don't panic when a pointer to a valid struct")

	assert.Panics(func() {
		s := time.Now()
		_ = Wrap(&s)
	}, "panic when not a valid struct")
}

func TestWrapper(t *testing.T) {
	assert := assert.New(t)

	loc, _ := time.LoadLocation("")

	res1 := &mockType1{
		ID:     "res123",
		Str:    "a_string",
		Int:    2,
		Int8:   8,
		Int16:  16,
		Int32:  32,
		Int64:  64,
		Uint:   4,
		Uint8:  8,
		Uint16: 16,
		Uint32: 32,
		Uint64: 64,
		Bool:   true,
		Time:   time.Date(2017, 1, 2, 3, 4, 5, 6, loc),
	}

	wrap1 := Wrap(res1)

	// ID and type
	id, typ := wrap1.IDAndType()
	assert.Equal(res1.ID, id, "id")
	assert.Equal("mocktypes1", typ, "type")

	wrap1.SetID("another-id")
	assert.Equal(res1.ID, "another-id", "set id")

	// Get attributes
	attr := wrap1.Attr("str")
	assert.Equal(Attr{
		Name:     "str",
		Type:     AttrTypeString,
		Nullable: false,
	}, attr, "get attribute (str)")
	assert.Equal(Attr{}, wrap1.Attr("nonexistent"), "get non-existent attribute")

	// Get relationships
	rel := wrap1.Rel("to-one")
	assert.Equal(Rel{
		FromName: "to-one",
		ToType:   "mocktypes2",
		ToOne:    true,
		ToName:   "",
		FromType: "mocktypes1",
		FromOne:  false,
	}, rel, "get relationship (to-one)")
	assert.Equal(Rel{}, wrap1.Rel("nonexistent"), "get non-existent relationship")

	// Get values (attributes)
	v1 := reflect.ValueOf(res1).Elem()
	for i := 0; i < v1.NumField(); i++ {
		f := v1.Field(i)
		sf := v1.Type().Field(i)
		n := sf.Tag.Get("json")

		if sf.Tag.Get("api") == "attr" {
			assert.Equal(f.Interface(), wrap1.Get(n), "api tag")
		}
	}

	// Set values (attributes)
	wrap1.Set("str", "another_string")
	assert.Equal("another_string", wrap1.Get("str"), "set string attribute")
	wrap1.Set("int", 3)
	assert.Equal(3, wrap1.Get("int"), "set int attribute")

	aStr := "another_string_ptr"
	aInt := int(123)
	aInt8 := int8(88)
	aInt16 := int16(1616)
	aInt32 := int32(3232)
	aInt64 := int64(6464)
	aUint := uint(456)
	aUint8 := uint8(88)
	aUint16 := uint16(1616)
	aUint32 := uint32(3232)
	aUint64 := uint64(6464)
	aBool := false
	aTime := time.Date(2018, 2, 3, 4, 5, 6, 7, loc)

	// Set the values (attributes) after the wrapping
	res2 := &mockType2{
		ID:        "res123",
		StrPtr:    &aStr,
		IntPtr:    &aInt,
		Int8Ptr:   &aInt8,
		Int16Ptr:  &aInt16,
		Int32Ptr:  &aInt32,
		Int64Ptr:  &aInt64,
		UintPtr:   &aUint,
		Uint8Ptr:  &aUint8,
		Uint16Ptr: &aUint16,
		Uint32Ptr: &aUint32,
		Uint64Ptr: &aUint64,
		BoolPtr:   &aBool,
		TimePtr:   &aTime,
	}

	wrap2 := Wrap(res2)

	// ID and type
	id, typ = wrap2.IDAndType()
	assert.Equal(res2.ID, id, "id 2")
	assert.Equal("mocktypes2", typ, "type 2")

	// Get values (attributes)
	v2 := reflect.ValueOf(res2).Elem()
	for i := 0; i < v2.NumField(); i++ {
		f := v2.Field(i)
		sf := v2.Type().Field(i)
		n := sf.Tag.Get("json")

		if sf.Tag.Get("api") == "attr" {
			assert.Equal(f.Interface(), wrap2.Get(n), "api tag 2")
		}
	}

	// Set values (attributes)
	var (
		anotherString = "anotherString"
		newInt        = 3
	)

	wrap2.Set("strptr", &anotherString)
	assert.Equal(&anotherString, wrap2.Get("strptr"), "set string pointer attribute")

	wrap2.Set("intptr", &newInt)
	assert.Equal(&newInt, wrap2.Get("intptr"), "set int pointer attribute")

	wrap2.Set("uintptr", nil)

	if wrap2.Get("uintptr") != nil {
		// We first do a != nil check because that's what we are really
		// checking and reflect.DeepEqual doesn't work exactly work the same
		// way. If the nil check fails, then the next line will fail too.
		assert.Equal("nil pointer", nil, wrap2.Get("uintptr"))
	}

	if res2.UintPtr != nil {
		// We first do a != nil check because that's what we are really
		// checking and reflect.DeepEqual doesn't work exactly work the same
		// way. If the nil check fails, then the next line will fail too.
		assert.Equal("nil pointer 2", nil, res2.UintPtr)
	}

	// New
	wrap3 := wrap1.New()

	for _, attr := range wrap1.Attrs() {
		assert.Equal(wrap1.Attr(attr.Name), wrap3.Attr(attr.Name), "copied attribute")
	}

	for _, rel := range wrap1.Rels() {
		assert.Equal(wrap1.Rel(rel.FromName), wrap3.Rel(rel.FromName), "copied relationship")
	}

	// Copy
	wrap3 = wrap1.Copy()

	for _, attr := range wrap1.Attrs() {
		assert.Equal(wrap1.Attr(attr.Name), wrap3.Attr(attr.Name), "copied attribute")
	}

	for _, rel := range wrap1.Rels() {
		assert.Equal(wrap1.Rel(rel.FromName), wrap3.Rel(rel.FromName), "copied relationship")
	}

	wrap3.Set("str", "another string")
	assert.NotEqual(
		wrap1.Get("str"),
		wrap3.Get("str"),
		fmt.Sprintf("modified value does not affect original"),
	)
}

func TestWrapperSet(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		typ   string // "1" for mockType1, "2" for mockType2
		field string
		val   interface{}
	}{
		{typ: "1", field: "str", val: "astring"},
		{typ: "1", field: "int", val: int(9)},
		{typ: "1", field: "int8", val: int8(9)},
		{typ: "1", field: "int16", val: int16(9)},
		{typ: "1", field: "int32", val: int32(9)},
		{typ: "1", field: "int64", val: int64(9)},
		{typ: "1", field: "uint", val: uint(9)},
		{typ: "1", field: "uint8", val: uint8(9)},
		{typ: "1", field: "uint16", val: uint16(9)},
		{typ: "1", field: "uint32", val: uint32(9)},
		{typ: "1", field: "uint64", val: uint64(9)},
		{typ: "1", field: "bool", val: bool(true)},
		{typ: "2", field: "strptr", val: ptr("astring")},
		{typ: "2", field: "intptr", val: ptr(int(9))},
		{typ: "2", field: "int8ptr", val: ptr(int8(9))},
		{typ: "2", field: "int16ptr", val: ptr(int16(9))},
		{typ: "2", field: "int32ptr", val: ptr(int32(9))},
		{typ: "2", field: "int64ptr", val: ptr(int64(9))},
		{typ: "2", field: "uintptr", val: ptr(uint(9))},
		{typ: "2", field: "uint8ptr", val: ptr(uint8(9))},
		{typ: "2", field: "uint16ptr", val: ptr(uint16(9))},
		{typ: "2", field: "uint32ptr", val: ptr(uint32(9))},
		{typ: "2", field: "uint64ptr", val: ptr(uint64(9))},
		{typ: "2", field: "boolptr", val: ptr(bool(true))},
	}

	for _, test := range tests {
		if test.typ == "1" {
			res1 := Wrap(&mockType1{})
			res1.Set(test.field, test.val)
			assert.EqualValues(test.val, res1.Get(test.field))
		}
	}
}

func TestWrapperGetAndSetErrors(t *testing.T) {
	assert := assert.New(t)

	mt := &mocktype{}
	wrap := Wrap(mt)

	// Get on empty field name
	assert.Panics(func() {
		_ = wrap.Get("")
	})

	// Get on unknown field name
	assert.Panics(func() {
		_ = wrap.Get("unknown")
	})

	// Set on empty field name
	assert.Panics(func() {
		wrap.Set("", "")
	})

	// Set on unknown field name
	assert.Panics(func() {
		wrap.Set("unknown", "")
	})

	// Set with value of wrong type
	assert.Panics(func() {
		wrap.Set("str", 42)
	})

	// GetToOne on empty field name
	assert.Panics(func() {
		_ = wrap.GetToOne("")
	})

	// GetToOne on unknown field name
	assert.Panics(func() {
		_ = wrap.GetToOne("unknown")
	})

	// GetToOne on attribute
	assert.Panics(func() {
		_ = wrap.GetToOne("str")
	})

	// GetToOne on to-many relationship
	assert.Panics(func() {
		_ = wrap.GetToOne("to-x")
	})

	// GetToMany on empty field name
	assert.Panics(func() {
		_ = wrap.GetToMany("")
	})

	// GetToMany on unknown field name
	assert.Panics(func() {
		_ = wrap.GetToMany("unknown")
	})

	// GetToMany on attribute
	assert.Panics(func() {
		_ = wrap.GetToMany("str")
	})

	// GetToMany on to-one relationship
	assert.Panics(func() {
		_ = wrap.GetToMany("to-1")
	})

	// SetToOne on empty field name
	assert.Panics(func() {
		wrap.SetToOne("", "id")
	})

	// SetToOne on unknown field name
	assert.Panics(func() {
		wrap.SetToOne("unknown", "id")
	})

	// SetToOne on attribute
	assert.Panics(func() {
		wrap.SetToOne("str", "id")
	})

	// SetToOne on to-many relationship
	assert.Panics(func() {
		wrap.SetToOne("to-x", "id")
	})

	// SetToMany on empty field name
	assert.Panics(func() {
		wrap.SetToMany("", []string{"id"})
	})

	// SetToMany on unknown field name
	assert.Panics(func() {
		wrap.SetToMany("unknown", []string{"id"})
	})

	// SetToMany on attribute
	assert.Panics(func() {
		wrap.SetToMany("str", []string{"id"})
	})

	// SetToMany on to-one relationship
	assert.Panics(func() {
		wrap.SetToMany("to-1", []string{"id"})
	})
}
