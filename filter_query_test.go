package jsonapi_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestFilterResource(t *testing.T) {
	assert := assert.New(t)

	runs := map[string]struct {
		kinds        []string
		vals         []string
		ops          []string
		val          string
		expectations []int
	}{
		"string": {
			kinds:        []string{"string"},
			vals:         []string{"aaa", "bbb", "bbb", "ccc", "ccc", "ccc"},
			ops:          []string{"=", "!=", "<", "<=", ">", ">="},
			val:          "bbb",
			expectations: []int{2, 4, 1, 3, 3, 5},
		},
		"int": {
			kinds:        []string{"int", "int8", "int16", "int32", "int64"},
			vals:         []string{"-1", "-1", "0", "1", "1", "2", "2", "3", "4"},
			ops:          []string{"=", "!=", "<", "<=", ">", ">="},
			val:          "1",
			expectations: []int{2, 7, 3, 5, 4, 6},
		},
		"uint": {
			kinds:        []string{"uint", "uint8", "uint16", "uint32", "uint64"},
			vals:         []string{"0", "1", "1", "2", "2", "3", "4"},
			ops:          []string{"=", "!=", "<", "<=", ">", ">="},
			val:          "1",
			expectations: []int{2, 5, 1, 3, 4, 6},
		},
		"bool": {
			kinds:        []string{"bool"},
			vals:         []string{"false", "true", "true"},
			ops:          []string{"=", "!="},
			val:          "true",
			expectations: []int{2, 1},
		},
		"time.Time": {
			kinds: []string{"time.Time"},
			vals: []string{
				"2009-11-10 22:59:58 +0000 UTC",
				"2009-11-10 22:59:59 +0000 UTC",
				"2009-11-10 23:00:00 +0000 UTC",
				"2009-11-10 23:00:01 +0000 UTC",
				"2009-11-10 23:00:01 +0000 UTC",
				"2009-11-10 23:00:02 +0000 UTC",
			},
			ops:          []string{"=", "!=", "<", "<=", ">", ">="},
			val:          "2009-11-10 23:00:00 +0000 UTC",
			expectations: []int{1, 5, 2, 3, 3, 4},
		},
		// "*string": {
		// 	kinds:        []string{"*string"},
		// 	vals:         []string{nil, "", "aaa", "bbb", "bbb", "ccc", "ccc", "ccc"},
		// 	ops:          []string{"=", "!=", "<", "<=", ">", ">="},
		// 	val:          "bbb",
		// 	expectations: []int{2, 4, 1, 3, 3, 5},
		// },
		// "*int": {
		// 	kinds:        []string{"*int", "*int8", "*int16", "*int32", "*int64"},
		// 	vals:         []string{nil, "", "-1", "-1", "0", "1", "1", "2", "2", "3", "4"},
		// 	ops:          []string{"=", "!=", "<", "<=", ">", ">="},
		// 	val:          "1",
		// 	expectations: []int{2, 7, 3, 5, 4, 6},
		// },
		// "*uint": {
		// 	kinds:        []string{"*uint", "*uint8", "*uint16", "*uint32", "*uint64"},
		// 	vals:         []string{nil, "", "0", "1", "1", "2", "2", "3", "4"},
		// 	ops:          []string{"=", "!=", "<", "<=", ">", ">="},
		// 	val:          "1",
		// 	expectations: []int{2, 5, 1, 3, 4, 6},
		// },
		// "*bool": {
		// 	kinds:        []string{"*bool"},
		// 	vals:         []string{nil, "", "false", "true", "true"},
		// 	ops:          []string{"=", "!="},
		// 	val:          "true",
		// 	expectations: []int{2, 1},
		// },
		// "*time.Time": {
		// 	kinds: []string{"*time.Time"},
		// 	vals: []string{
		// 		nil
		// 		"2009-11-10 22:59:58 +0000 UTC",
		// 		"2009-11-10 22:59:59 +0000 UTC",
		// 		"2009-11-10 23:00:00 +0000 UTC",
		// 		"2009-11-10 23:00:01 +0000 UTC",
		// 		"2009-11-10 23:00:01 +0000 UTC",
		// 		"2009-11-10 23:00:02 +0000 UTC",
		// 	},
		// 	ops:          []string{"=", "!=", "<", "<=", ">", ">="},
		// 	val:          "2009-11-10 23:00:00 +0000 UTC",
		// 	expectations: []int{1, 5, 2, 3, 3, 4},
		// },
	}

	for typ, run := range runs {
		for _, kind := range run.kinds {
			src := &SoftCollection{}
			src.Type = &Type{Name: "type"}
			ty, n := GetAttrType(typ)
			src.Type.Attrs = map[string]Attr{
				"attr": Attr{
					Name: "attr",
					Type: ty,
					Null: n,
				},
			}

			for _, v := range run.vals {
				res := &SoftResource{}
				res.SetType(src.Type)
				res.Set("attr", makeVal(v, kind))
				src.Add(res)
			}

			cond := &Condition{}
			cond.Field = "attr"
			cond.Val = makeVal(run.val, kind)
			for i, op := range run.ops {
				cond.Op = op
				dst := &SoftCollection{}
				for i := 0; i < src.Len(); i++ {
					r := src.Elem(i)
					if FilterResource(r, cond) {
						dst.Add(r.Copy())
					}
				}
				assert.Equal(
					run.expectations[i],
					dst.Len(),
					fmt.Sprintf("%s %v (%s)", cond.Op, cond.Val, kind),
				)
			}
		}
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

func makeVal(v string, t string) interface{} {
	var r interface{}
	switch t {
	// String
	case "string":
		r = v
	// Integers
	case "int":
		r, _ = strconv.Atoi(v)
	case "int8":
		r, _ = strconv.Atoi(v)
		r = int8(r.(int))
	case "int16":
		r, _ = strconv.Atoi(v)
		r = int16(r.(int))
	case "int32":
		r, _ = strconv.Atoi(v)
		r = int32(r.(int))
	case "int64":
		r, _ = strconv.Atoi(v)
		r = int64(r.(int))
	case "uint":
		r, _ = strconv.Atoi(v)
		r = uint(r.(int))
	case "uint8":
		r, _ = strconv.Atoi(v)
		r = uint8(r.(int))
	case "uint16":
		r, _ = strconv.Atoi(v)
		r = uint16(r.(int))
	case "uint32":
		r, _ = strconv.Atoi(v)
		r = uint32(r.(int))
	case "uint64":
		r, _ = strconv.Atoi(v)
		r = uint64(r.(int))
	// Bool
	case "bool":
		r = v == "true"
	// time.Time
	case "time.Time":
		r, _ = time.Parse("2006-02-01 15:04:05 -0700 MST", v)
	}
	return r
}
