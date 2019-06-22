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
		vals         []string
		ops          []string
		val          interface{}
		expectations []int
	}{
		"string": {
			vals:         []string{"aaa", "bbb", "bbb", "ccc", "ccc", "ccc"},
			ops:          []string{"=", "!=", "<", "<=", ">", ">="},
			val:          "bbb",
			expectations: []int{2, 4, 1, 3, 3, 5},
		},
		"int": {
			vals:         []string{"-1", "-1", "0", "1", "1", "2", "2", "3", "4"},
			ops:          []string{"=", "!=", "<", "<=", ">", ">="},
			val:          1,
			expectations: []int{2, 4, 1, 3, 3, 5},
		},
	}

	for typ, run := range runs {
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

		// strings := []string{"aaa", "bbb", "bbb", "ccc", "ccc", "ccc"}
		// ops := []string{"=", "!=", "<", "<=", ">", ">="}
		expectations := []int{2, 4, 1, 3, 3, 5}

		for _, v := range run.vals {
			res := &SoftResource{}
			res.SetType(src.Type)
			res.Set("attr", makeVal(v, typ))
			src.Add(res)
		}

		cond := &Condition{}
		cond.Field = "attr"
		cond.Val = run.val
		for i, op := range run.ops {
			cond.Op = op
			dst := &SoftCollection{}
			for i := 0; i < src.Len(); i++ {
				r := src.Elem(i)
				if FilterResource(r, cond) {
					dst.Add(r.Copy())
				}
			}
			assert.Equal(expectations[i], dst.Len(), fmt.Sprintf("%s %v (string)", cond.Op, cond.Val))
		}
	}

	// Int
	// intKinds := []string{"int", "int8", "int16", "int32", "int64"}
	// for _, intKind := range intKinds {
	// 	src = &SoftCollection{}
	// 	src.Type = &Type{Name: "type"}
	// 	typ, n = GetAttrType(intKind)
	// 	src.Type.Attrs = map[string]Attr{
	// 		"attr": Attr{
	// 			Name: "attr",
	// 			Type: typ,
	// 			Null: n,
	// 		},
	// 	}

	// 	ints := []int{-1, -1, 0, 1, 1, 2, 2, 3, 4}
	// 	ops = []string{"=", "!=", "<", "<=", ">", ">="}
	// 	val = makeVal(1, intKind)
	// 	expectations := []int{2, 7, 3, 5, 4, 6}

	// 	for _, v := range ints {
	// 		res := &SoftResource{}
	// 		res.SetType(src.Type)
	// 		res.Set("attr", makeVal(v, intKind))
	// 		src.Add(res)
	// 	}

	// 	cond = &Condition{}
	// 	cond.Field = "attr"
	// 	cond.Val = val
	// 	for i := range ops {
	// 		cond.Op = ops[i]
	// 		dst := &SoftCollection{}
	// 		for i := 0; i < src.Len(); i++ {
	// 			r := src.Elem(i)
	// 			if FilterResource(r, cond) {
	// 				dst.Add(r.Copy())
	// 			}
	// 		}
	// 		assert.EqualValues(
	// 			expectations[i],
	// 			dst.Len(),
	// 			fmt.Sprintf("%s %v (%s)", cond.Op, cond.Val, intKind),
	// 		)
	// 	}
	// }

	// // Uint
	// uintKinds := []string{"uint", "uint8", "uint16", "uint32", "uint64"}
	// for _, uintKind := range uintKinds {
	// 	src = &SoftCollection{}
	// 	src.Type = &Type{Name: "type"}
	// 	typ, n = GetAttrType(uintKind)
	// 	src.Type.Attrs = map[string]Attr{
	// 		"attr": Attr{
	// 			Name: "attr",
	// 			Type: typ,
	// 			Null: n,
	// 		},
	// 	}

	// 	uints := []uint{0, 1, 1, 2, 2, 3, 4}
	// 	ops = []string{"=", "!=", "<", "<=", ">", ">="}
	// 	val = makeVal(1, uintKind)
	// 	expectations := []uint{2, 5, 1, 3, 4, 6}

	// 	for _, v := range uints {
	// 		res := &SoftResource{}
	// 		res.SetType(src.Type)
	// 		res.Set("attr", makeVal(int(v), uintKind))
	// 		src.Add(res)
	// 	}

	// 	cond = &Condition{}
	// 	cond.Field = "attr"
	// 	cond.Val = val
	// 	for i := range ops {
	// 		cond.Op = ops[i]
	// 		dst := &SoftCollection{}
	// 		for i := 0; i < src.Len(); i++ {
	// 			r := src.Elem(i)
	// 			if FilterResource(r, cond) {
	// 				dst.Add(r.Copy())
	// 			}
	// 		}
	// 		assert.EqualValues(
	// 			expectations[i],
	// 			dst.Len(),
	// 			fmt.Sprintf("%s %v (%s)", cond.Op, cond.Val, uintKind),
	// 		)
	// 	}
	// }

	// // Bool
	// src = &SoftCollection{}
	// src.Type = &Type{Name: "type"}
	// typ, n = GetAttrType("bool")
	// src.Type.Attrs = map[string]Attr{
	// 	"attr": Attr{
	// 		Name: "attr",
	// 		Type: typ,
	// 		Null: n,
	// 	},
	// }

	// bools := []bool{false, true, true}
	// ops = []string{"=", "!="}
	// val = true
	// expectations = []int{2, 1}

	// for _, v := range bools {
	// 	res := &SoftResource{}
	// 	res.SetType(src.Type)
	// 	res.Set("attr", v)
	// 	src.Add(res)
	// }

	// cond = &Condition{}
	// cond.Field = "attr"
	// cond.Val = val
	// for i := range ops {
	// 	cond.Op = ops[i]
	// 	dst := &SoftCollection{}
	// 	for i := 0; i < src.Len(); i++ {
	// 		r := src.Elem(i)
	// 		if FilterResource(r, cond) {
	// 			dst.Add(r.Copy())
	// 		}
	// 	}
	// 	assert.Equal(expectations[i], dst.Len(), fmt.Sprintf("%s %v (bool)", cond.Op, cond.Val))
	// }

	// // Time
	// src = &SoftCollection{}
	// src.Type = &Type{Name: "type"}
	// typ, n = GetAttrType("bool")
	// src.Type.Attrs = map[string]Attr{
	// 	"attr": Attr{
	// 		Name: "attr",
	// 		Type: typ,
	// 		Null: n,
	// 	},
	// }

	// now := time.Now()
	// times := []time.Time{
	// 	now.Add(-2 * time.Second),
	// 	now.Add(-time.Second),
	// 	now,
	// 	now.Add(time.Second),
	// 	now.Add(time.Second),
	// 	now.Add(2 * time.Second),
	// }
	// ops = []string{"=", "!=", "<", "<=", ">", ">="}
	// val = now
	// expectations = []int{1, 5, 2, 3, 3, 4}

	// for _, v := range times {
	// 	res := &SoftResource{}
	// 	res.SetType(src.Type)
	// 	res.Set("attr", v)
	// 	src.Add(res)
	// }

	// cond = &Condition{}
	// cond.Field = "attr"
	// cond.Val = val
	// for i := range ops {
	// 	cond.Op = ops[i]
	// 	dst := &SoftCollection{}
	// 	for i := 0; i < src.Len(); i++ {
	// 		r := src.Elem(i)
	// 		if FilterResource(r, cond) {
	// 			dst.Add(r.Copy())
	// 		}
	// 	}
	// 	assert.Equal(expectations[i], dst.Len(), fmt.Sprintf("%s %v (time.Time)", cond.Op, cond.Val))
	// }
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
		r, _ = time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", v)
	}
	return r
}
