package jsonapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/kkaribu/tchek"
)

func TestFilterQuery(t *testing.T) {
	// time1, _ := time.Parse(time.RFC3339Nano, "2012-05-16T17:45:28.2539Z")
	// time2, _ := time.Parse(time.RFC3339Nano, "2013-06-24T22:03:34.8276Z")

	tests := []struct {
		query             string
		expectedCondition Condition
		expectedError     bool
	}{
		{
			// 0
			query:         ``,
			expectedError: true,
		},
		{
			// 1
			query:         `{"v":null}`,
			expectedError: false, // TODO
		},
		{
			// 2
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

	for n, test := range tests {
		cdt := Condition{}
		err := json.Unmarshal([]byte(test.query), &cdt)

		tchek.ErrorExpected(t, n, test.expectedError, err)

		if !test.expectedError {
			tchek.AreEqual(t, n, test.expectedCondition, cdt)

			data, err := json.Marshal(&cdt)
			tchek.UnintendedError(err)

			tchek.AreEqual(t, n, tchek.MakeOneLineNoSpaces(test.query), tchek.MakeOneLineNoSpaces(string(data)))
		}
	}

	// Test marshaling error
	_, err := json.Marshal(&Condition{
		Op:  "=",
		Val: func() {},
	})
	tchek.ErrorExpected(t, -1, true, err)

	_, err = json.Marshal(&Condition{
		Op:  "",
		Val: "",
	})
	tchek.ErrorExpected(t, -2, false, err) // TODO
}

func BenchmarkMarshalFilterQuery(b *testing.B) {
	cdt := Condition{
		Op: "or",
		Val: []Condition{
			Condition{
				Op:  "in",
				Val: []string{"a", "b", "c"},
			},
			Condition{
				Op: "and",
				Val: []Condition{
					Condition{
						Op:  "~",
						Val: "%a",
					},
					Condition{
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
