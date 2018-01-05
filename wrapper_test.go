package jsonapi

import (
	"reflect"
	"testing"
	"time"

	"github.com/kkaribu/tchek"
)

func TestWrapper(t *testing.T) {
	loc, _ := time.LoadLocation("")

	res1 := &MockType1{
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
	tchek.AreEqual(t, -1, res1.ID, id)
	tchek.AreEqual(t, -1, "mocktypes1", typ)

	// Get attributes
	v1 := reflect.ValueOf(res1).Elem()
	for i := 0; i < v1.NumField(); i++ {
		f := v1.Field(i)
		sf := v1.Type().Field(i)
		n := sf.Tag.Get("json")

		if sf.Tag.Get("api") == "attr" {
			tchek.AreEqual(t, -1, f.Interface(), wrap1.Get(n))
		}
	}

	// Set attributes
	wrap1.Set("str", "another_string")
	tchek.AreEqual(t, -1, "another_string", wrap1.Get("str"))
	wrap1.Set("int", 3)
	tchek.AreEqual(t, -1, 3, wrap1.Get("int"))

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
	res2 := &MockType2{
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
	tchek.AreEqual(t, -1, res2.ID, id)
	tchek.AreEqual(t, -1, "mocktypes2", typ)

	// Get attributes
	v2 := reflect.ValueOf(res2).Elem()
	for i := 0; i < v2.NumField(); i++ {
		f := v2.Field(i)
		sf := v2.Type().Field(i)
		n := sf.Tag.Get("json")

		if sf.Tag.Get("api") == "attr" {
			tchek.AreEqual(t, 2, f.Interface(), wrap2.Get(n))
		}
	}

	// Set attributes
	var anotherString = "anotherString"
	wrap2.Set("strptr", &anotherString)
	tchek.AreEqual(t, -1, &anotherString, wrap2.Get("strptr"))
	var newInt = 3
	wrap2.Set("intptr", &newInt)
	tchek.AreEqual(t, -1, &newInt, wrap2.Get("intptr"))

	// Copy
	wrap3 := wrap1.Copy()

	for _, attr := range wrap1.Attrs() {
		tchek.AreEqual(t, -1, wrap1.Get(attr.Name), wrap3.Get(attr.Name))

		if attr.Type == "bool" {
			wrap3.Set(attr.Name, !wrap1.Get(attr.Name).(bool))
		} else if attr.Type == "*bool" {
			wrap3.Set(attr.Name, !*(wrap1.Get(attr.Name).(*bool)))
		} else if attr.Type == "time.Time" || attr.Type == "*time.Time" {
			wrap3.Set(attr.Name, time.Now())
		} else {
			wrap3.Set(attr.Name, "0")
		}
		tchek.AreNotEqual(t, -1, wrap1.Get(attr.Name), wrap3.Get(attr.Name))
	}
}
