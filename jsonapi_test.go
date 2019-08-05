package jsonapi_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestMarshaling(t *testing.T) {
	tests := []struct {
		name   string
		schema *Schema
		col    Collection
		url    string
		// doc         *Document
		expected    string
		expectedErr error
	}{
		{
			name:   "all zero",
			schema: getEmptyBaseSchema(),
			col:    getEmptyBaseCollection(),
			url:    "/type1",
			// doc:         &Document{},
			expected:    "expected",
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run("aaa", func(t *testing.T) {
			assert := assert.New(t)
			_ = fmt.Sprintf("test: %+v\n", test)

			url, err := NewURLFromRaw(test.schema, test.url)
			assert.NoError(err)

			doc := &Document{}

			// Marshaling
			payload, err := Marshal(doc, url)
			assert.Equal(err, test.expectedErr)
			assert.Equal(string(payload), test.expected)

			// Unmarshaling
			// doc2, err := Unmarshal(payload, url, test.schema)
			// assert.NoError(err)
			// assert.Equal(doc, doc2)
		})
	}
}

func getEmptyBaseSchema() *Schema {
	schema := &Schema{}
	return schema
}

func getBaseSchema() *Schema {
	schema := &Schema{}
	return schema
}

func getEmptyBaseCollection() Collection {
	col := &SoftCollection{}
	return col
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
