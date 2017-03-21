package jsonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/kkaribu/tchek"
)

func TestMarshalResource(t *testing.T) {
	loc, _ := time.LoadLocation("")
	reg := NewMockRegistry()

	tests := []struct {
		data          Resource
		host          string
		params        string
		meta          map[string]interface{}
		jsonapi       map[string]interface{}
		errorExpected bool
		payloadFile   string
	}{
		{
			// 0
			data: mocktypes1.Elem(0),
			meta: map[string]interface{}{
				"num":       42,
				"timestamp": time.Date(2017, 1, 2, 3, 4, 5, 6, loc),
				"tf":        true,
				"str":       "a string",
			},
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payloadFile:   "resource-1",
		}, {
			// 1
			data: mocktypes2.Elem(1),
			host: "https://example.org",
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payloadFile:   "resource-2",
		}, {
			// 2
			data:   mocktypes2.Elem(1),
			host:   "https://example.org",
			params: "?fields[mocktypes2]=strptr,uintptr,int",
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payloadFile:   "resource-3",
		},
	}

	for n, test := range tests {
		id, resType := test.data.IDAndType()
		rawurl := fmt.Sprintf("/%s/%s%s", resType, id, test.params)

		url, err := ParseRawURL(reg, rawurl)
		tchek.UnintendedError(err)

		// Options
		opts := NewOptions(test.host, url.Params)
		opts.Meta = test.meta
		opts.JSONAPI = test.jsonapi

		// Marshal
		payload, err := Marshal(test.data, url, opts)
		tchek.ErrorExpected(t, n, test.errorExpected, err)

		if !test.errorExpected {
			var out bytes.Buffer

			// Format the payload
			json.Indent(&out, payload, "", "\t")
			output := out.String()

			// Retrieve the expected result from file
			content, err := ioutil.ReadFile("tests/" + test.payloadFile + ".json")
			tchek.UnintendedError(err)
			out.Reset()
			json.Indent(&out, content, "", "\t")
			// Trim because otherwise there is an extra line at the end
			expectedOutput := strings.TrimSpace(out.String())

			tchek.AreEqual(t, n, expectedOutput, output)
		}
	}
}

func TestMarshalCollection(t *testing.T) {
	loc, _ := time.LoadLocation("")
	reg := NewMockRegistry()

	tests := []struct {
		data          Collection
		host          string
		params        string
		meta          map[string]interface{}
		jsonapi       map[string]interface{}
		errorExpected bool
		payloadFile   string
	}{
		{
			// 0
			data: mocktypes1,
			meta: map[string]interface{}{
				"num":       -32820,
				"timestamp": time.Date(1981, 2, 3, 4, 5, 6, 0, loc),
				"tf":        false,
				"str":       "//\n\téç.\\",
			},
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payloadFile:   "collection-1",
		}, {
			// 1
			data:   mocktypes2,
			host:   "https://example.org",
			params: "?fields[mocktypes2]=uintptr,boolptr,timeptr",
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payloadFile:   "collection-2",
		}, {
			// 2
			data: WrapCollection(Wrap(&MockType1{})),
			host: "https://example.org",
			jsonapi: map[string]interface{}{
				"version": "1.0",
			},
			errorExpected: false,
			payloadFile:   "collection-3",
		},
	}

	for n, test := range tests {
		_, resType := test.data.Sample().IDAndType()
		rawurl := fmt.Sprintf("/%s%s", resType, test.params)

		url, err := ParseRawURL(reg, rawurl)
		tchek.UnintendedError(err)

		// Options
		opts := NewOptions(test.host, url.Params)
		opts.Meta = test.meta
		opts.JSONAPI = test.jsonapi

		// Marshal
		payload, err := Marshal(test.data, url, opts)
		tchek.ErrorExpected(t, n, test.errorExpected, err)

		if !test.errorExpected {
			var out bytes.Buffer

			// Format the payload
			json.Indent(&out, payload, "", "\t")
			output := out.String()

			// Retrieve the expected result from file
			content, err := ioutil.ReadFile("tests/" + test.payloadFile + ".json")
			tchek.UnintendedError(err)
			out.Reset()
			json.Indent(&out, content, "", "\t")
			// Trim because otherwise there is an extra line at the end
			expectedOutput := strings.TrimSpace(out.String())

			tchek.AreEqual(t, n, expectedOutput, output)
		}
	}
}
