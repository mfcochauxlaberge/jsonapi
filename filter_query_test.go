package jsonapi_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestFilterResource(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()

	tests := []struct {
		rval     interface{}
		op       string
		cval     interface{}
		expected bool
	}{
		// string
		{rval: "bbb", op: "=", cval: "aaa", expected: false},
		{rval: "bbb", op: "=", cval: "bbb", expected: true},
		{rval: "bbb", op: "!=", cval: "aaa", expected: true},
		{rval: "bbb", op: "!=", cval: "bbb", expected: false},
		{rval: "bbb", op: "<", cval: "aaa", expected: false},
		{rval: "bbb", op: "<", cval: "bbb", expected: false},
		{rval: "bbb", op: "<", cval: "ccc", expected: true},
		{rval: "bbb", op: "<=", cval: "aaa", expected: false},
		{rval: "bbb", op: "<=", cval: "bbb", expected: true},
		{rval: "bbb", op: "<=", cval: "ccc", expected: true},
		{rval: "bbb", op: ">", cval: "aaa", expected: true},
		{rval: "bbb", op: ">", cval: "bbb", expected: false},
		{rval: "bbb", op: ">", cval: "ccc", expected: false},
		{rval: "bbb", op: ">=", cval: "aaa", expected: true},
		{rval: "bbb", op: ">=", cval: "bbb", expected: true},
		{rval: "bbb", op: ">=", cval: "ccc", expected: false},

		// int
		{rval: 1, op: "=", cval: 0, expected: false},
		{rval: 1, op: "=", cval: 1, expected: true},
		{rval: 1, op: "!=", cval: 0, expected: true},
		{rval: 1, op: "!=", cval: 1, expected: false},
		{rval: 1, op: "<", cval: 0, expected: false},
		{rval: 1, op: "<", cval: 1, expected: false},
		{rval: 1, op: "<", cval: 3, expected: true},
		{rval: 1, op: "<=", cval: 0, expected: false},
		{rval: 1, op: "<=", cval: 1, expected: true},
		{rval: 1, op: "<=", cval: 3, expected: true},
		{rval: 1, op: ">", cval: 0, expected: true},
		{rval: 1, op: ">", cval: 1, expected: false},
		{rval: 1, op: ">", cval: 3, expected: false},
		{rval: 1, op: ">=", cval: 0, expected: true},
		{rval: 1, op: ">=", cval: 1, expected: true},
		{rval: 1, op: ">=", cval: 3, expected: false},

		// int8
		{rval: int8(1), op: "=", cval: int8(0), expected: false},
		{rval: int8(1), op: "=", cval: int8(1), expected: true},
		{rval: int8(1), op: "!=", cval: int8(0), expected: true},
		{rval: int8(1), op: "!=", cval: int8(1), expected: false},
		{rval: int8(1), op: "<", cval: int8(0), expected: false},
		{rval: int8(1), op: "<", cval: int8(1), expected: false},
		{rval: int8(1), op: "<", cval: int8(2), expected: true},
		{rval: int8(1), op: "<=", cval: int8(0), expected: false},
		{rval: int8(1), op: "<=", cval: int8(1), expected: true},
		{rval: int8(1), op: "<=", cval: int8(2), expected: true},
		{rval: int8(1), op: ">", cval: int8(0), expected: true},
		{rval: int8(1), op: ">", cval: int8(1), expected: false},
		{rval: int8(1), op: ">", cval: int8(2), expected: false},
		{rval: int8(1), op: ">=", cval: int8(0), expected: true},
		{rval: int8(1), op: ">=", cval: int8(1), expected: true},
		{rval: int8(1), op: ">=", cval: int8(2), expected: false},

		// int16
		{rval: int16(1), op: "=", cval: int16(0), expected: false},
		{rval: int16(1), op: "=", cval: int16(1), expected: true},
		{rval: int16(1), op: "!=", cval: int16(0), expected: true},
		{rval: int16(1), op: "!=", cval: int16(1), expected: false},
		{rval: int16(1), op: "<", cval: int16(0), expected: false},
		{rval: int16(1), op: "<", cval: int16(1), expected: false},
		{rval: int16(1), op: "<", cval: int16(2), expected: true},
		{rval: int16(1), op: "<=", cval: int16(0), expected: false},
		{rval: int16(1), op: "<=", cval: int16(1), expected: true},
		{rval: int16(1), op: "<=", cval: int16(2), expected: true},
		{rval: int16(1), op: ">", cval: int16(0), expected: true},
		{rval: int16(1), op: ">", cval: int16(1), expected: false},
		{rval: int16(1), op: ">", cval: int16(2), expected: false},
		{rval: int16(1), op: ">=", cval: int16(0), expected: true},
		{rval: int16(1), op: ">=", cval: int16(1), expected: true},
		{rval: int16(1), op: ">=", cval: int16(2), expected: false},

		// int32
		{rval: int32(1), op: "=", cval: int32(0), expected: false},
		{rval: int32(1), op: "=", cval: int32(1), expected: true},
		{rval: int32(1), op: "!=", cval: int32(0), expected: true},
		{rval: int32(1), op: "!=", cval: int32(1), expected: false},
		{rval: int32(1), op: "<", cval: int32(0), expected: false},
		{rval: int32(1), op: "<", cval: int32(1), expected: false},
		{rval: int32(1), op: "<", cval: int32(2), expected: true},
		{rval: int32(1), op: "<=", cval: int32(0), expected: false},
		{rval: int32(1), op: "<=", cval: int32(1), expected: true},
		{rval: int32(1), op: "<=", cval: int32(2), expected: true},
		{rval: int32(1), op: ">", cval: int32(0), expected: true},
		{rval: int32(1), op: ">", cval: int32(1), expected: false},
		{rval: int32(1), op: ">", cval: int32(2), expected: false},
		{rval: int32(1), op: ">=", cval: int32(0), expected: true},
		{rval: int32(1), op: ">=", cval: int32(1), expected: true},
		{rval: int32(1), op: ">=", cval: int32(2), expected: false},

		// int64
		{rval: int64(1), op: "=", cval: int64(0), expected: false},
		{rval: int64(1), op: "=", cval: int64(1), expected: true},
		{rval: int64(1), op: "!=", cval: int64(0), expected: true},
		{rval: int64(1), op: "!=", cval: int64(1), expected: false},
		{rval: int64(1), op: "<", cval: int64(0), expected: false},
		{rval: int64(1), op: "<", cval: int64(1), expected: false},
		{rval: int64(1), op: "<", cval: int64(2), expected: true},
		{rval: int64(1), op: "<=", cval: int64(0), expected: false},
		{rval: int64(1), op: "<=", cval: int64(1), expected: true},
		{rval: int64(1), op: "<=", cval: int64(2), expected: true},
		{rval: int64(1), op: ">", cval: int64(0), expected: true},
		{rval: int64(1), op: ">", cval: int64(1), expected: false},
		{rval: int64(1), op: ">", cval: int64(2), expected: false},
		{rval: int64(1), op: ">=", cval: int64(0), expected: true},
		{rval: int64(1), op: ">=", cval: int64(1), expected: true},
		{rval: int64(1), op: ">=", cval: int64(2), expected: false},

		// uint
		{rval: uint(1), op: "=", cval: uint(0), expected: false},
		{rval: uint(1), op: "=", cval: uint(1), expected: true},
		{rval: uint(1), op: "!=", cval: uint(0), expected: true},
		{rval: uint(1), op: "!=", cval: uint(1), expected: false},
		{rval: uint(1), op: "<", cval: uint(0), expected: false},
		{rval: uint(1), op: "<", cval: uint(1), expected: false},
		{rval: uint(1), op: "<", cval: uint(2), expected: true},
		{rval: uint(1), op: "<=", cval: uint(0), expected: false},
		{rval: uint(1), op: "<=", cval: uint(1), expected: true},
		{rval: uint(1), op: "<=", cval: uint(2), expected: true},
		{rval: uint(1), op: ">", cval: uint(0), expected: true},
		{rval: uint(1), op: ">", cval: uint(1), expected: false},
		{rval: uint(1), op: ">", cval: uint(2), expected: false},
		{rval: uint(1), op: ">=", cval: uint(0), expected: true},
		{rval: uint(1), op: ">=", cval: uint(1), expected: true},
		{rval: uint(1), op: ">=", cval: uint(2), expected: false},

		// uint8
		{rval: uint8(1), op: "=", cval: uint8(0), expected: false},
		{rval: uint8(1), op: "=", cval: uint8(1), expected: true},
		{rval: uint8(1), op: "!=", cval: uint8(0), expected: true},
		{rval: uint8(1), op: "!=", cval: uint8(1), expected: false},
		{rval: uint8(1), op: "<", cval: uint8(0), expected: false},
		{rval: uint8(1), op: "<", cval: uint8(1), expected: false},
		{rval: uint8(1), op: "<", cval: uint8(2), expected: true},
		{rval: uint8(1), op: "<=", cval: uint8(0), expected: false},
		{rval: uint8(1), op: "<=", cval: uint8(1), expected: true},
		{rval: uint8(1), op: "<=", cval: uint8(2), expected: true},
		{rval: uint8(1), op: ">", cval: uint8(0), expected: true},
		{rval: uint8(1), op: ">", cval: uint8(1), expected: false},
		{rval: uint8(1), op: ">", cval: uint8(2), expected: false},
		{rval: uint8(1), op: ">=", cval: uint8(0), expected: true},
		{rval: uint8(1), op: ">=", cval: uint8(1), expected: true},
		{rval: uint8(1), op: ">=", cval: uint8(2), expected: false},

		// uint16
		{rval: uint16(1), op: "=", cval: uint16(0), expected: false},
		{rval: uint16(1), op: "=", cval: uint16(1), expected: true},
		{rval: uint16(1), op: "!=", cval: uint16(0), expected: true},
		{rval: uint16(1), op: "!=", cval: uint16(1), expected: false},
		{rval: uint16(1), op: "<", cval: uint16(0), expected: false},
		{rval: uint16(1), op: "<", cval: uint16(1), expected: false},
		{rval: uint16(1), op: "<", cval: uint16(2), expected: true},
		{rval: uint16(1), op: "<=", cval: uint16(0), expected: false},
		{rval: uint16(1), op: "<=", cval: uint16(1), expected: true},
		{rval: uint16(1), op: "<=", cval: uint16(2), expected: true},
		{rval: uint16(1), op: ">", cval: uint16(0), expected: true},
		{rval: uint16(1), op: ">", cval: uint16(1), expected: false},
		{rval: uint16(1), op: ">", cval: uint16(2), expected: false},
		{rval: uint16(1), op: ">=", cval: uint16(0), expected: true},
		{rval: uint16(1), op: ">=", cval: uint16(1), expected: true},
		{rval: uint16(1), op: ">=", cval: uint16(2), expected: false},

		// uint32
		{rval: uint32(1), op: "=", cval: uint32(0), expected: false},
		{rval: uint32(1), op: "=", cval: uint32(1), expected: true},
		{rval: uint32(1), op: "!=", cval: uint32(0), expected: true},
		{rval: uint32(1), op: "!=", cval: uint32(1), expected: false},
		{rval: uint32(1), op: "<", cval: uint32(0), expected: false},
		{rval: uint32(1), op: "<", cval: uint32(1), expected: false},
		{rval: uint32(1), op: "<", cval: uint32(2), expected: true},
		{rval: uint32(1), op: "<=", cval: uint32(0), expected: false},
		{rval: uint32(1), op: "<=", cval: uint32(1), expected: true},
		{rval: uint32(1), op: "<=", cval: uint32(2), expected: true},
		{rval: uint32(1), op: ">", cval: uint32(0), expected: true},
		{rval: uint32(1), op: ">", cval: uint32(1), expected: false},
		{rval: uint32(1), op: ">", cval: uint32(2), expected: false},
		{rval: uint32(1), op: ">=", cval: uint32(0), expected: true},
		{rval: uint32(1), op: ">=", cval: uint32(1), expected: true},
		{rval: uint32(1), op: ">=", cval: uint32(2), expected: false},

		// uint64
		{rval: uint64(1), op: "=", cval: uint64(0), expected: false},
		{rval: uint64(1), op: "=", cval: uint64(1), expected: true},
		{rval: uint64(1), op: "!=", cval: uint64(0), expected: true},
		{rval: uint64(1), op: "!=", cval: uint64(1), expected: false},
		{rval: uint64(1), op: "<", cval: uint64(0), expected: false},
		{rval: uint64(1), op: "<", cval: uint64(1), expected: false},
		{rval: uint64(1), op: "<", cval: uint64(2), expected: true},
		{rval: uint64(1), op: "<=", cval: uint64(0), expected: false},
		{rval: uint64(1), op: "<=", cval: uint64(1), expected: true},
		{rval: uint64(1), op: "<=", cval: uint64(2), expected: true},
		{rval: uint64(1), op: ">", cval: uint64(0), expected: true},
		{rval: uint64(1), op: ">", cval: uint64(1), expected: false},
		{rval: uint64(1), op: ">", cval: uint64(2), expected: false},
		{rval: uint64(1), op: ">=", cval: uint64(0), expected: true},
		{rval: uint64(1), op: ">=", cval: uint64(1), expected: true},
		{rval: uint64(1), op: ">=", cval: uint64(2), expected: false},

		// bool
		{rval: true, op: "=", cval: true, expected: true},
		{rval: true, op: "=", cval: false, expected: false},
		{rval: true, op: "!=", cval: true, expected: false},
		{rval: true, op: "!=", cval: false, expected: true},

		// time.Time
		{rval: now, op: "=", cval: now.Add(-time.Second), expected: false},
		{rval: now, op: "=", cval: now, expected: true},
		{rval: now, op: "!=", cval: now.Add(-time.Second), expected: true},
		{rval: now, op: "!=", cval: now, expected: false},
		{rval: now, op: "<", cval: now.Add(-time.Second), expected: false},
		{rval: now, op: "<", cval: now, expected: false},
		{rval: now, op: "<", cval: now.Add(time.Second), expected: true},
		{rval: now, op: "<=", cval: now.Add(-time.Second), expected: false},
		{rval: now, op: "<=", cval: now, expected: true},
		{rval: now, op: "<=", cval: now.Add(time.Second), expected: true},
		{rval: now, op: ">", cval: now.Add(-time.Second), expected: true},
		{rval: now, op: ">", cval: now, expected: false},
		{rval: now, op: ">", cval: now.Add(time.Second), expected: false},
		{rval: now, op: ">=", cval: now.Add(-time.Second), expected: true},
		{rval: now, op: ">=", cval: now, expected: true},
		{rval: now, op: ">=", cval: now.Add(time.Second), expected: false},

		// *string
		{rval: ptr("bbb"), op: "=", cval: nilptr("string"), expected: false},
		{rval: ptr("bbb"), op: "=", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: "=", cval: ptr("bbb"), expected: true},
		{rval: ptr("bbb"), op: "!=", cval: nilptr("string"), expected: true},
		{rval: ptr("bbb"), op: "!=", cval: ptr("aaa"), expected: true},
		{rval: ptr("bbb"), op: "!=", cval: ptr("bbb"), expected: false},
		{rval: ptr("bbb"), op: "<", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: "<", cval: ptr("bbb"), expected: false},
		{rval: ptr("bbb"), op: "<", cval: ptr("ccc"), expected: true},
		{rval: ptr("bbb"), op: "<=", cval: ptr("aaa"), expected: false},
		{rval: ptr("bbb"), op: "<=", cval: ptr("bbb"), expected: true},
		{rval: ptr("bbb"), op: "<=", cval: ptr("ccc"), expected: true},
		{rval: ptr("bbb"), op: ">", cval: ptr("aaa"), expected: true},
		{rval: ptr("bbb"), op: ">", cval: ptr("bbb"), expected: false},
		{rval: ptr("bbb"), op: ">", cval: ptr("ccc"), expected: false},
		{rval: ptr("bbb"), op: ">=", cval: ptr("aaa"), expected: true},
		{rval: ptr("bbb"), op: ">=", cval: ptr("bbb"), expected: true},
		{rval: ptr("bbb"), op: ">=", cval: ptr("ccc"), expected: false},
	}

	for _, test := range tests {
		typ := &Type{Name: "type"}
		ty, n := GetAttrType(fmt.Sprintf("%T", test.rval))
		typ.Attrs = map[string]Attr{
			"attr": Attr{
				Name: "attr",
				Type: ty,
				Null: n,
			},
		}

		res := &SoftResource{}
		res.SetType(typ)
		res.Set("attr", test.rval)

		cond := &Condition{}
		cond.Field = "attr"
		cond.Op = test.op
		cond.Val = test.cval

		assert.Equal(
			test.expected,
			FilterResource(res, cond),
			fmt.Sprintf("%v %s %v is %v", test.rval, test.op, test.cval, test.expected),
		)
	}
}

func TestFilterQuery(t *testing.T) {
	assert := assert.New(t)

	// time1, _ := time.Parse(time.RFC3339Nano, "2012-05-16T17:45:28.2539Z")
	// time2, _ := time.Parse(time.RFC3339Nano, "2013-06-24T22:03:34.8276Z")

	tests := []struct {
		name              string
		query             string
		expectedCondition Condition
		expectedError     bool
	}{
		{
			name:          "empty",
			query:         ``,
			expectedError: true,
		},
		{
			name:          "null value",
			query:         `{"v":null}`,
			expectedError: false, // TODO
		},
		{
			name: "standard values",
			query: `{
				"c": "col",
				"f": "field",
				"o": "=",
				"v": "string"
			}`,
			expectedCondition: Condition{
				Field: "field",
				Op:    "=",
				Val:   "string",
				Col:   "col",
			},
			expectedError: false,
		},
	}

	for _, test := range tests {
		cdt := Condition{}
		err := json.Unmarshal([]byte(test.query), &cdt)

		assert.Equal(test.expectedError, err != nil, test.name)

		if !test.expectedError {
			assert.Equal(test.expectedCondition, cdt, test.name)

			data, err := json.Marshal(&cdt)
			assert.NoError(err, test.name)

			assert.Equal(makeOneLineNoSpaces(test.query), makeOneLineNoSpaces(string(data)), test.name)
		}
	}

	// Test marshaling error
	_, err := json.Marshal(&Condition{
		Op:  "=",
		Val: func() {},
	})
	assert.Equal(true, err != nil, "function as value")

	_, err = json.Marshal(&Condition{
		Op:  "",
		Val: "",
	})
	assert.Equal(false, err != nil, "empty operation and value") // TODO
}

func BenchmarkMarshalFilterQuery(b *testing.B) {
	cdt := Condition{
		Op: "or",
		Val: []Condition{
			{
				Op:  "in",
				Val: []string{"a", "b", "c"},
			},
			{
				Op: "and",
				Val: []Condition{
					{
						Op:  "~",
						Val: "%a",
					},
					{
						Op:  ">=",
						Val: "u",
					},
				},
			},
		},
	}

	var (
		data []byte
		err  error
	)

	for n := 0; n < b.N; n++ {
		data, err = json.Marshal(cdt)
	}

	fmt.Fprintf(ioutil.Discard, "%v %v", data, err)
}

func BenchmarkUnmarshalFilterQuery(b *testing.B) {
	query := []byte(`
		{ "or": [
			{ "in": ["a", "b", "c"] },
			{ "and": [
				{ "~": "%a" },
				{ "\u003e=": "u" }
			] }
		] }
	`)

	var (
		cdt Condition
		err error
	)

	for n := 0; n < b.N; n++ {
		cdt = Condition{}
		err = json.Unmarshal(query, &cdt)
	}

	fmt.Fprintf(ioutil.Discard, "%v %v", cdt, err)
}

func ptr(v interface{}) interface{} {
	switch c := v.(type) {
	// String
	case string:
		return &c
	// Integers
	case int:
		return &c
	case int8:
		return &c
	case int16:
		return &c
	case int32:
		return &c
	case int64:
		return &c
	case uint:
		return &c
	case uint8:
		return &c
	case uint16:
		return &c
	case uint32:
		return &c
	case uint64:
		return &c
	// Bool
	case bool:
		return &c
	// time.Time
	case time.Time:
		return &c
	}
	return nil
}

func nilptr(t string) interface{} {
	switch t {
	// String
	case "string":
		var p *string
		return p
	// Integers
	case "int":
		var p *int
		return p
	case "int8":
		var p *int8
		return p
	case "int16":
		var p *int16
		return p
	case "int32":
		var p *int32
		return p
	case "int64":
		var p *int64
		return p
	case "uint":
		var p *uint
		return p
	case "uint8":
		var p *uint8
		return p
	case "uint16":
		var p *uint16
		return p
	case "uint32":
		var p *uint32
		return p
	case "uint64":
		var p *uint64
		return p
	// Bool
	case "bool":
		var p *bool
		return p
	// time.Time
	case "time.Time":
		var p *time.Time
		return p
	}
	return nil
}
