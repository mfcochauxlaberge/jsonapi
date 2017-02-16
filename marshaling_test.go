package jsonapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"kkaribu/tchek"
)

func TestMarshal(t *testing.T) {
	loc, _ := time.LoadLocation("")

	tests := []struct {
		src           interface{}
		params        *Params
		meta          map[string]interface{}
		jsonapi       map[string]interface{}
		errorExpected bool
		payload       string
	}{
		{
			// 0
			src:    users.Elem(0),
			params: params[0],
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payload:       "payload-0",
		}, {
			// 1
			src:    users,
			params: params[1],
			meta: map[string]interface{}{
				"b":            true,
				"generated-at": time.Date(2017, 01, 02, 03, 04, 05, 0, loc),
				"str":          "str",
				"number":       42,
			},
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payload:       "payload-1",
		}, {
			// 2
			src: Identifier{
				Type: "users",
				ID:   "1",
			},
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payload:       "payload-2",
		}, {
			// 3
			src: NewIdentifiers("users", []string{"1", "a", "abc123"}),
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payload:       "payload-3",
		}, {
			// 4
			src: Identifier{},
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payload:       "payload-4",
		}, {
			// 5
			src: Identifiers{},
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payload:       "payload-5",
		},
	}

	for n, test := range tests {
		// Marshal
		payload, err := Marshal(test.src, test.params, Extra{
			Meta: test.meta,
			JSONAPI: map[string]interface{}{
				"version": "1.0",
			},
		})
		tchek.ErrorExpected(t, n, test.errorExpected, err)

		if !test.errorExpected {
			var out bytes.Buffer

			// Format the payload
			json.Indent(&out, payload, "", "\t")
			output := out.String()

			// Retrieve the expected result from file
			content, err := ioutil.ReadFile("tests/" + test.payload + ".json")
			tchek.UnintendedError(err)
			out.Reset()
			json.Indent(&out, content, "", "\t")
			// Trim because otherwise there is an extra line at the end
			expectedOutput := strings.TrimSpace(out.String())

			tchek.AreEqual(t, n, expectedOutput, output)
		}
	}
}
