package jsonapi_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

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
		Bool:   true,
		Time:   time.Date(2017, 1, 2, 3, 4, 5, 6, loc),
	}

	wrap1 := Wrap(res1)

	// ID and type
	id, typ := wrap1.IDAndType()
	assert.Equal(res1.ID, id, "id")
	assert.Equal("mocktypes1", typ, "type")

	// Get attributes
	v1 := reflect.ValueOf(res1).Elem()
	for i := 0; i < v1.NumField(); i++ {
		f := v1.Field(i)
		sf := v1.Type().Field(i)
		n := sf.Tag.Get("json")

		if sf.Tag.Get("api") == "attr" {
			assert.Equal(f.Interface(), wrap1.Get(n), "api tag")
		}
	}

	// Set attributes
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
	aBool := false
	aTime := time.Date(2018, 2, 3, 4, 5, 6, 7, loc)

	// Set the attributes after the wrapping
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
		BoolPtr:   &aBool,
		TimePtr:   &aTime,
	}

	wrap2 := Wrap(res2)

	// ID and type
	id, typ = wrap2.IDAndType()
	assert.Equal(res2.ID, id, "id 2")
	assert.Equal("mocktypes2", typ, "type 2")

	// Get attributes
	v2 := reflect.ValueOf(res2).Elem()
	for i := 0; i < v2.NumField(); i++ {
		f := v2.Field(i)
		sf := v2.Type().Field(i)
		n := sf.Tag.Get("json")

		if sf.Tag.Get("api") == "attr" {
			assert.Equal(f.Interface(), wrap2.Get(n), "api tag 2")
		}
	}

	// Set attributes
	var anotherString = "anotherString"
	wrap2.Set("strptr", &anotherString)
	assert.Equal(&anotherString, wrap2.Get("strptr"), "set string pointer attribute")
	var newInt = 3
	wrap2.Set("intptr", &newInt)
	assert.Equal(&newInt, wrap2.Get("intptr"), "set int pointer attribute")
	wrap2.Set("uintptr", nil)
	if wrap2.Get("uintptr") != nil {
		// We first do a != nil check because that's what we are really
		// checking and reflect.DeepEqual doesn't work exactly work the same
		// way. If the nil check fails, then the next line will fail too.
		assert.Equal(t, "nil pointer", nil, wrap2.Get("uintptr"))
	}
	if res2.UintPtr != nil {
		// We first do a != nil check because that's what we are really
		// checking and reflect.DeepEqual doesn't work exactly work the same
		// way. If the nil check fails, then the next line will fail too.
		assert.Equal(t, "nil pointer 2", nil, res2.UintPtr)
	}

	// Copy
	wrap3 := wrap1.Copy()

	for _, attr := range wrap1.Attrs() {
		assert.Equal(wrap1.Get(attr.Name), wrap3.Get(attr.Name), "copied attribute")

		if attr.Type == AttrTypeBool && !attr.Null {
			wrap3.Set(attr.Name, !wrap1.Get(attr.Name).(bool))
		} else if attr.Type == AttrTypeBool && attr.Null {
			wrap3.Set(attr.Name, !*(wrap1.Get(attr.Name).(*bool)))
		} else if attr.Type == AttrTypeTime {
			wrap3.Set(attr.Name, time.Now())
		} else {
			wrap3.Set(attr.Name, "0")
		}
		assert.NotEqual(wrap1.Get(attr.Name), wrap3.Get(attr.Name), fmt.Sprintf("modified copied attribute %s (%v)", attr.Name, attr.Type))
	}
}
