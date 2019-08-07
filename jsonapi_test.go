package jsonapi_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestMarshaling(t *testing.T) {
	update := false

	paths := []string{
		"/type1",
		"/type1/t1-1",
		"/type1/t1-1/to-1",
		"/type1/t1-1/relationships/to-1",
	}

	tests := []struct {
		name   string
		schema *Schema
		col    *SoftCollection
		url    string
	}{}

	for _, path := range paths {
		tests = append(tests, struct {
			name   string
			schema *Schema
			col    *SoftCollection
			url    string
		}{
			name:   "some name",
			schema: getBaseSchema(),
			col:    getEmptyBaseCollection("type1"),
			url:    path,
		})
	}

	for i := range tests {
		i := i
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// URL
			url, err := NewURLFromRaw(test.schema, test.url)
			assert.NoError(err)

			// Data
			var data interface{}
			if url.IsCol {
				// If it's a collection
				page := test.col.Range(nil, nil, nil, nil, 1000, 0)
				dataCol := &SoftCollection{
					Type: test.col.Type,
				}
				for i := range page {
					dataCol.Add(page[i])
				}
				data = dataCol
			} else {
				// If it's a resource
				for i := 0; i < test.col.Len(); i++ {
					if test.col.At(i).GetID() == url.ResID {
						data = test.col.At(i)
						break
					}
				}
			}

			// Document
			doc := &Document{
				Data: data,
			}

			// Marshaling
			payload, err := Marshal(doc, url)
			assert.NoError(err)

			// Golden file
			filename := "test" + strconv.Itoa(i) + ".json"
			path := filepath.Join("testdata", "goldenfiles", filename)
			if !update {
				// Retrieve the expected result from file
				expected, _ := ioutil.ReadFile(path)
				assert.NoError(err, test.name)
				assert.JSONEq(string(expected), string(payload))
			} else {
				dst := &bytes.Buffer{}
				err = json.Indent(dst, payload, "", "\t")
				assert.NoError(err)
				// TODO Figure out whether 0644 is okay or not.
				err = ioutil.WriteFile(path, dst.Bytes(), 0644)
				assert.NoError(err)
			}

			// Unmarshaling
			// TODO
			// doc2, err := Unmarshal(payload, url, test.schema)
			// assert.NoError(err)
			// assert.Equal(doc, doc2)
		})
	}
}

func getBaseSchema() *Schema {
	schema := &Schema{}
	_ = schema.AddType(MustReflect(type1{}))
	_ = schema.AddType(MustReflect(type2{}))
	_ = schema.AddType(MustReflect(type3{}))
	if len(schema.Check()) > 0 {
		panic("base schema for tests should be valid")
	}
	return schema
}

func getEmptyBaseCollection(t string) *SoftCollection {
	schema := getBaseSchema()
	typ := schema.GetType(t)
	col := &SoftCollection{
		Type: &typ,
	}
	return col
}

// type1 is a fake struct that defines a JSON:API type for test purposes.
type type1 struct {
	ID string `json:"id" api:"type1"`

	// Attributes
	Str    string    `json:"str" api:"attr"`
	Int    int       `json:"int" api:"attr"`
	Int8   int8      `json:"int8" api:"attr"`
	Int16  int16     `json:"int16" api:"attr"`
	Int32  int32     `json:"int32" api:"attr"`
	Int64  int64     `json:"int64" api:"attr"`
	Uint   uint      `json:"uint" api:"attr"`
	Uint8  uint8     `json:"uint8" api:"attr"`
	Uint16 uint16    `json:"uint16" api:"attr"`
	Uint32 uint32    `json:"uint32" api:"attr"`
	Uint64 uint64    `json:"uint64" api:"attr"`
	Bool   bool      `json:"bool" api:"attr"`
	Time   time.Time `json:"time" api:"attr"`

	// Relationships
	To1      string   `json:"to-1" api:"rel,type2"`
	To1From1 string   `json:"to-1-from-1" api:"rel,type2,to-1-from-1"`
	To1FromX string   `json:"to-1-from-x" api:"rel,type2,to-x-from-1"`
	ToX      []string `json:"to-x" api:"rel,type2"`
	ToXFrom1 []string `json:"to-x-from-1" api:"rel,type2,to-1-from-x"`
	ToXFromX []string `json:"to-x-from-x" api:"rel,type2,to-x-from-x"`
}

// type2 is a fake struct that defines a JSON:API type for test purposes.
type type2 struct {
	ID string `json:"id" api:"type2"`

	// Attributes
	StrPtr    *string    `json:"strptr" api:"attr"`
	IntPtr    *int       `json:"intptr" api:"attr"`
	Int8Ptr   *int8      `json:"int8ptr" api:"attr"`
	Int16Ptr  *int16     `json:"int16ptr" api:"attr"`
	Int32Ptr  *int32     `json:"int32ptr" api:"attr"`
	Int64Ptr  *int64     `json:"int64ptr" api:"attr"`
	UintPtr   *uint      `json:"uintptr" api:"attr"`
	Uint8Ptr  *uint8     `json:"uint8ptr" api:"attr"`
	Uint16Ptr *uint16    `json:"uint16ptr" api:"attr"`
	Uint32Ptr *uint32    `json:"uint32ptr" api:"attr"`
	Uint64Ptr *uint64    `json:"uint64ptr" api:"attr"`
	BoolPtr   *bool      `json:"boolptr" api:"attr"`
	TimePtr   *time.Time `json:"timeptr" api:"attr"`

	// Relationships
	To1From1 string   `json:"to-1-from-1" api:"rel,type1,to-1-from-1"`
	To1FromX string   `json:"to-1-from-x" api:"rel,type1,to-x-from-1"`
	ToXFrom1 []string `json:"to-x-from-1" api:"rel,type1,to-1-from-x"`
	ToXFromX []string `json:"to-x-from-x" api:"rel,type1,to-x-from-x"`
}

// type3 is a fake struct that defines a JSON:API type for test purposes.
type type3 struct {
	ID string `json:"id" api:"type3"`

	// Attributes
	Attr1 string `json:"attr1" api:"attr"`
	Attr2 int    `json:"attr2" api:"attr"`

	// Relationships
	Rel1 string   `json:"rel1" api:"rel,type1"`
	Rel2 []string `json:"rel2" api:"rel,type1"`
}

var (
	mocktypes1  Collection
	mocktypes11 Collection
	mocktypes2  Collection
	mocktypes21 Collection
	mocktypes3  Collection
	// urls  []*URL
)

func init() {
	loc, _ := time.LoadLocation("")

	// Resources
	mocktypes1 = WrapCollection(Wrap(&mockType1{}))
	mocktypes1.Add(
		Wrap(&mockType1{
			ID: "mt1-1",
			// Use default (zero) value for each attribute
		}),
	)
	mocktypes1.Add(
		Wrap(&mockType1{
			ID:     "mt1-2",
			Str:    "",
			Int:    -42,
			Int8:   80,
			Int16:  160,
			Int32:  320,
			Int64:  6464640000,
			Uint:   42,
			Uint8:  8,
			Uint16: 1600,
			Uint32: 32000,
			Uint64: 64000,
			Bool:   false,
			Time:   time.Date(2017, 1, 2, 3, 4, 5, 6, loc),
		}),
	)

	mocktypes2 = WrapCollection(Wrap(&mockType2{}))
	mocktypes2.Add(
		Wrap(&mockType2{
			ID: "mt2-1",
			// Use nil values
		}),
	)
	strPtr := "str"
	intPtr := int(-42)
	int8Ptr := int8(80)
	int16Ptr := int16(160)
	int32Ptr := int32(320)
	int64Ptr := int64(6464640000)
	uintPtr := uint(42)
	uint8Ptr := uint8(8)
	uint16Ptr := uint16(1600)
	uint32Ptr := uint32(32000)
	uint64Ptr := uint64(64000)
	boolPtr := false
	timePtr := time.Date(2017, 1, 2, 3, 4, 5, 6, loc)
	mocktypes2.Add(
		Wrap(&mockType2{
			ID:        "mt1-2",
			StrPtr:    &strPtr,
			IntPtr:    &intPtr,
			Int8Ptr:   &int8Ptr,
			Int16Ptr:  &int16Ptr,
			Int32Ptr:  &int32Ptr,
			Int64Ptr:  &int64Ptr,
			UintPtr:   &uintPtr,
			Uint8Ptr:  &uint8Ptr,
			Uint16Ptr: &uint16Ptr,
			Uint32Ptr: &uint32Ptr,
			Uint64Ptr: &uint64Ptr,
			BoolPtr:   &boolPtr,
			TimePtr:   &timePtr,
		}),
	)

	mocktypes3 = WrapCollection(Wrap(&mockType3{}))
	mocktypes3.Add(
		Wrap(&mockType3{
			ID: "mt3-1",
		}),
	)
	mocktypes3.Add(
		Wrap(&mockType3{
			ID:    "mt3-1",
			Attr1: "str",
			Attr2: 32,
		}),
	)

	mocktypes11 = WrapCollection(Wrap(&mockType1{}))
	mocktypes11.Add(
		Wrap(&mockType1{
			ID:             "mt11-1",
			ToOne:          "mt21-1",
			ToMany:         []string{"mt21-2", "mt21-3"},
			ToOneFromOne:   "",
			ToManyFromMany: []string{},
			ToManyFromOne:  []string{},
		}),
	)

	mocktypes21 = WrapCollection(Wrap(&mockType2{}))
	mocktypes21.Add(
		Wrap(&mockType2{
			ID:      "mt21-1",
			BoolPtr: &boolPtr,
		}),
	)
	mocktypes21.Add(
		Wrap(&mockType2{
			ID:     "mt21-2",
			IntPtr: &intPtr,
		}),
	)
	mocktypes21.Add(
		Wrap(&mockType2{
			ID:     "mt21-3",
			StrPtr: &strPtr,
		}),
	)
}
