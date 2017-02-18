package jsonapi

import (
	"reflect"
	"testing"
	"time"

	"kkaribu/tchek"
)

func TestWrapper(t *testing.T) {
	loc, _ := time.LoadLocation("")

	r := &res{}
	w := Wrap(r)

	aStr := "a_string_ptr"
	aInt := int(2)
	aInt8 := int8(8)
	aInt16 := int16(16)
	aInt32 := int32(32)
	aInt64 := int64(64)
	aUint := uint(4)
	aUint8 := uint8(8)
	aUint16 := uint16(16)
	aUint32 := uint32(32)
	aBool := true
	aTime := time.Date(2017, 1, 2, 3, 4, 5, 6, loc)

	// Set the attributes after the wrapping
	r.ID = "res123"
	r.Str = "a_string"
	r.StrPtr = &aStr
	r.Int = 2
	r.Int8 = 8
	r.Int16 = 16
	r.Int32 = 32
	r.Int64 = 64
	r.IntPtr = &aInt
	r.Int8Ptr = &aInt8
	r.Int16Ptr = &aInt16
	r.Int32Ptr = &aInt32
	r.Int64Ptr = &aInt64
	r.Uint = 4
	r.Uint8 = 8
	r.Uint16 = 16
	r.Uint32 = 32
	r.Uint64 = 64
	r.UintPtr = &aUint
	r.Uint8Ptr = &aUint8
	r.Uint16Ptr = &aUint16
	r.Uint32Ptr = &aUint32
	r.Bool = true
	r.BoolPtr = &aBool
	r.Time = time.Date(2017, 1, 2, 3, 4, 5, 6, loc)
	r.TimePtr = &aTime

	// ID and type
	id, typ := w.IDAndType()
	tchek.AreEqual(t, 0, r.ID, id)
	tchek.AreEqual(t, 0, "res", typ)

	v := reflect.ValueOf(r).Elem()

	// Attributes
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		sf := v.Type().Field(i)
		n := sf.Tag.Get("json")

		if sf.Tag.Get("api") == "attr" {
			tchek.AreEqual(t, 0, f.Interface(), w.Get(n))
		}
	}
}

type res struct {
	ID string `json:"id" api:"res"`

	Str    string  `json:"str" api:"attr"`
	StrPtr *string `json:"strptr" api:"attr"`

	Int      int    `json:"int" api:"attr"`
	Int8     int8   `json:"int8" api:"attr"`
	Int16    int16  `json:"int16" api:"attr"`
	Int32    int32  `json:"int32" api:"attr"`
	Int64    int64  `json:"int64" api:"attr"`
	IntPtr   *int   `json:"intptr" api:"attr"`
	Int8Ptr  *int8  `json:"int8ptr" api:"attr"`
	Int16Ptr *int16 `json:"int16ptr" api:"attr"`
	Int32Ptr *int32 `json:"int32ptr" api:"attr"`
	Int64Ptr *int64 `json:"int64ptr" api:"attr"`

	Uint      uint    `json:"uint" api:"attr"`
	Uint8     uint8   `json:"uint8" api:"attr"`
	Uint16    uint16  `json:"uint16" api:"attr"`
	Uint32    uint32  `json:"uint32" api:"attr"`
	Uint64    uint64  `json:"uint64" api:"attr"`
	UintPtr   *uint   `json:"uintptr" api:"attr"`
	Uint8Ptr  *uint8  `json:"uint8ptr" api:"attr"`
	Uint16Ptr *uint16 `json:"uint16ptr" api:"attr"`
	Uint32Ptr *uint32 `json:"uint32ptr" api:"attr"`

	Bool    bool  `json:"bool" api:"attr"`
	BoolPtr *bool `json:"boolptr" api:"attr"`

	Time    time.Time  `json:"time" api:"attr"`
	TimePtr *time.Time `json:"timeptr" api:"attr"`

	ToOne  string   `json:"toone" api:"rel,target,inverse"`
	ToMany []string `json:"tomany" api:"rel,target,inverse"`
}
