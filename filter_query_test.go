package jsonapi_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"
	"github.com/mfcochauxlaberge/tchek"
)

func TestFilterQuery(t *testing.T) {
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

		tchek.ErrorExpected(t, test.name, test.expectedError, err)

		if !test.expectedError {
			tchek.AreEqual(t, test.name, test.expectedCondition, cdt)

			data, err := json.Marshal(&cdt)
			tchek.UnintendedError(err)

			tchek.AreEqual(t, test.name, tchek.MakeOneLineNoSpaces(test.query), tchek.MakeOneLineNoSpaces(string(data)))
		}
	}

	// Test marshaling error
	_, err := json.Marshal(&Condition{
		Op:  "=",
		Val: func() {},
	})
	tchek.ErrorExpected(t, "function as value", true, err)

	_, err = json.Marshal(&Condition{
		Op:  "",
		Val: "",
	})
	tchek.ErrorExpected(t, "empty operation and value", false, err) // TODO
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
