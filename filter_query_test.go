package jsonapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/kkaribu/tchek"
)

func TestFilterQuery(t *testing.T) {
	time1, _ := time.Parse(time.RFC3339Nano, "2012-05-16T17:45:28.2539Z")
	time2, _ := time.Parse(time.RFC3339Nano, "2013-06-24T22:03:34.8276Z")

	tests := []struct {
		query             string
		kind              string
		expectedCondition Condition
		expectedError     bool
	}{
		{
			// 0
			query: `{ "=": "a" }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   "=",
				Val:  "a",
			},
			expectedError: false,
		}, {
			// 1
			query: `{ "!=": "a" }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   "!=",
				Val:  "a",
			},
			expectedError: false,
		}, {
			// 2
			query: `{ "\u003c": "a" }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   "<",
				Val:  "a",
			},
			expectedError: false,
		}, {
			// 3
			query: `{ "\u003e": "a" }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   ">",
				Val:  "a",
			},
			expectedError: false,
		}, {
			// 4
			query: `{ "\u003c=": "a" }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   "<=",
				Val:  "a",
			},
			expectedError: false,
		}, {
			// 5
			query: `{ "\u003e=": "a" }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   ">=",
				Val:  "a",
			},
			expectedError: false,
		}, {
			// 6
			query: `{ "in": ["a", "b", "c"] }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   "in",
				Val:  []string{"a", "b", "c"},
			},
			expectedError: false,
		}, {
			// 7
			query: `{ "notin": ["a", "b", "c"] }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   "notin",
				Val:  []string{"a", "b", "c"},
			},
			expectedError: false,
		}, {
			// 8
			query: `{ "~": "a%" }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   "~",
				Val:  "a%",
			},
			expectedError: false,
		}, {
			// 9
			query: `{ "!~": "a%" }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   "!~",
				Val:  "a%",
			},
			expectedError: false,
		}, {
			// 10
			query: `{ "=": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "=",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 11
			query: `{ "!=": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "!=",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 12
			query: `{ "\u003c": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "<",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 13
			query: `{ "\u003e": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   ">",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 14
			query: `{ "\u003c=": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "<=",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 15
			query: `{ "\u003e=": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   ">=",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 16
			query: `{ "in": [1, 2, 3] }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "in",
				Val:  []int{1, 2, 3},
			},
			expectedError: false,
		}, {
			// 17
			query: `{ "notin": [1, 2, 3] }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "notin",
				Val:  []int{1, 2, 3},
			},
			expectedError: false,
		}, {
			// 18
			query: `{ "=": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "=",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 19
			query: `{ "!=": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "!=",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 20
			query: `{ "\u003c": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "<",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 21
			query: `{ "\u003e": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   ">",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 22
			query: `{ "\u003c=": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "<=",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 23
			query: `{ "\u003e=": 1 }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   ">=",
				Val:  1,
			},
			expectedError: false,
		}, {
			// 24
			query: `{ "in": [1, 2, 3] }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "in",
				Val:  []int{1, 2, 3},
			},
			expectedError: false,
		}, {
			// 25
			query: `{ "notin": [1, 2, 3] }`,
			kind:  "number",
			expectedCondition: Condition{
				Kind: "number",
				Op:   "notin",
				Val:  []int{1, 2, 3},
			},
			expectedError: false,
		}, {
			// 26
			query: `{ "=": "2012-05-16T17:45:28.2539Z" }`,
			kind:  "time",
			expectedCondition: Condition{
				Kind: "time",
				Op:   "=",
				Val:  time1,
			},
			expectedError: false,
		}, {
			// 27
			query: `{ "!=": "2012-05-16T17:45:28.2539Z" }`,
			kind:  "time",
			expectedCondition: Condition{
				Kind: "time",
				Op:   "!=",
				Val:  time1,
			},
			expectedError: false,
		}, {
			// 28
			query: `{ "\u003c": "2012-05-16T17:45:28.2539Z" }`,
			kind:  "time",
			expectedCondition: Condition{
				Kind: "time",
				Op:   "<",
				Val:  time1,
			},
			expectedError: false,
		}, {
			// 29
			query: `{ "\u003e": "2012-05-16T17:45:28.2539Z" }`,
			kind:  "time",
			expectedCondition: Condition{
				Kind: "time",
				Op:   ">",
				Val:  time1,
			},
			expectedError: false,
		}, {
			// 30
			query: `{ "\u003c=": "2012-05-16T17:45:28.2539Z" }`,
			kind:  "time",
			expectedCondition: Condition{
				Kind: "time",
				Op:   "<=",
				Val:  time1,
			},
			expectedError: false,
		}, {
			// 31
			query: `{ "\u003e=": "2012-05-16T17:45:28.2539Z" }`,
			kind:  "time",
			expectedCondition: Condition{
				Kind: "time",
				Op:   ">=",
				Val:  time1,
			},
			expectedError: false,
		}, {
			// 32
			query: `{ "in": ["2012-05-16T17:45:28.2539Z", "2013-06-24T22:03:34.8276Z"] }`,
			kind:  "time",
			expectedCondition: Condition{
				Kind: "time",
				Op:   "in",
				Val:  []time.Time{time1, time2},
			},
			expectedError: false,
		}, {
			// 33
			query: `{ "notin": ["2012-05-16T17:45:28.2539Z", "2013-06-24T22:03:34.8276Z"] }`,
			kind:  "time",
			expectedCondition: Condition{
				Kind: "time",
				Op:   "notin",
				Val:  []time.Time{time1, time2},
			},
			expectedError: false,
		}, {
			// 34
			query: `{ "or": [
				{ "in": ["a", "b", "c"] },
				{ "and": [
					{ "~": "%a" },
					{ "\u003e=": "u" }
				] }
			] }`,
			kind: "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   "or",
				Val: []Condition{
					Condition{
						Kind: "string",
						Op:   "in",
						Val:  []string{"a", "b", "c"},
					},
					Condition{
						Kind: "string",
						Op:   "and",
						Val: []Condition{
							Condition{
								Kind: "string",
								Op:   "~",
								Val:  "%a",
							},
							Condition{
								Kind: "string",
								Op:   ">=",
								Val:  "u",
							},
						},
					},
				},
			},
			expectedError: false,
		}, {
			// 35
			query:         `.abcdez`,
			expectedError: true,
		}, {
			// 36
			query:         `[]`,
			expectedError: true,
		}, {
			// 37
			query:         `{}`,
			expectedError: true,
		}, {
			// 38
			query:         `{ "and": "a" }`,
			expectedError: true,
		}, {
			// 39
			query:         `{ "and": ["a"] }`,
			expectedError: true,
		}, {
			// 40
			query:         `{ "=": ["a"] }`,
			kind:          "string",
			expectedError: true,
		}, {
			// 41
			query:         `{ "in": "a" }`,
			kind:          "string",
			expectedError: true,
		}, {
			// 42
			query:         `{ "?": "a" }`,
			kind:          "string",
			expectedError: true,
		}, {
			// 43
			query:         `{ "=": "a" }`,
			kind:          "number",
			expectedError: true,
		}, {
			// 44
			query:         `{ "in": "a" }`,
			kind:          "number",
			expectedError: true,
		}, {
			// 45
			query:         `{ "?": "a" }`,
			kind:          "number",
			expectedError: true,
		}, {
			// 46
			query:         `{ "=": "a" }`,
			kind:          "time",
			expectedError: true,
		}, {
			// 47
			query:         `{ "in": "a" }`,
			kind:          "time",
			expectedError: true,
		}, {
			// 48
			query:         `{ "?": "a" }`,
			kind:          "time",
			expectedError: true,
		}, {
			// 49
			query:         `{ "=": "a" }`,
			kind:          "?",
			expectedError: true,
		}, {
			// 50
			query:         ``,
			kind:          "string",
			expectedError: true,
		}, {
			// 51
			query:         ``,
			kind:          "?",
			expectedError: true,
		}, {
			// 52
			query: `{ "=": null }`,
			kind:  "string",
			expectedCondition: Condition{
				Kind: "string",
				Op:   "=",
				Val:  nil,
			},
			expectedError: false,
		},
	}

	for n, test := range tests {
		cdt := Condition{Kind: test.kind}
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
		Kind: "string",
		Op:   "=",
		Val:  func() {},
	})
	tchek.ErrorExpected(t, -1, true, err)

	_, err = json.Marshal(&Condition{
		Kind: "string",
		Op:   "",
		Val:  "",
	})
	tchek.ErrorExpected(t, -1, true, err)
}

func BenchmarkMarshalFilterQuery(b *testing.B) {
	cdt := Condition{
		Kind: "string",
		Op:   "or",
		Val: []Condition{
			Condition{
				Kind: "string",
				Op:   "in",
				Val:  []string{"a", "b", "c"},
			},
			Condition{
				Kind: "string",
				Op:   "and",
				Val: []Condition{
					Condition{
						Kind: "string",
						Op:   "~",
						Val:  "%a",
					},
					Condition{
						Kind: "string",
						Op:   ">=",
						Val:  "u",
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
		cdt = Condition{Kind: "string"}
		err = json.Unmarshal(query, &cdt)
	}

	fmt.Fprintf(ioutil.Discard, "%v %v", cdt, err)
}
