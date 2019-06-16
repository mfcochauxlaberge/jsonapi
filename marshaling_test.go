package jsonapi_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"
	"github.com/mfcochauxlaberge/tchek"
)

func TestMarshalResource(t *testing.T) {
	loc, _ := time.LoadLocation("")
	schema := newMockSchema()

	tests := []struct {
		name          string
		data          Resource
		prepath       string
		params        string
		meta          map[string]interface{}
		errorExpected bool
		payloadFile   string
	}{
		{
			name: "resource with meta",
			data: mocktypes1.Elem(0),
			meta: map[string]interface{}{
				"num":       42,
				"timestamp": time.Date(2017, 1, 2, 3, 4, 5, 6, loc),
				"tf":        true,
				"str":       "a string",
			},
			errorExpected: false,
			payloadFile:   "resource-1",
		}, {
			name:          "resource with prepath",
			data:          mocktypes2.Elem(1),
			prepath:       "https://example.org",
			errorExpected: false,
			payloadFile:   "resource-2",
		}, {
			name:          "resource with prepath and params",
			data:          mocktypes2.Elem(1),
			prepath:       "https://example.org",
			params:        "?fields[mocktypes2]=strptr,uintptr,int",
			errorExpected: false,
			payloadFile:   "resource-3",
		},
	}

	for _, test := range tests {
		doc := NewDocument()
		doc.PrePath = test.prepath

		doc.Data = test.data

		id := test.data.GetID()
		resType := test.data.GetType()
		rawurl := fmt.Sprintf("%s/%s/%s%s", test.prepath, resType, id, test.params)

		url, err := ParseRawURL(schema, rawurl)
		tchek.UnintendedError(err)

		doc.Meta = test.meta

		// Marshal
		payload, err := Marshal(doc, url)
		tchek.ErrorExpected(t, test.name, test.errorExpected, err)

		if !test.errorExpected {
			var out bytes.Buffer

			// Format the payload
			json.Indent(&out, payload, "", "\t")
			output := out.String()

			// Retrieve the expected result from file
			content, err := ioutil.ReadFile("testdata/" + test.payloadFile + ".json")
			tchek.UnintendedError(err)
			out.Reset()
			json.Indent(&out, content, "", "\t")
			// Trim because otherwise there is an extra line at the end
			expectedOutput := strings.TrimSpace(out.String())

			tchek.AreEqual(t, test.name, expectedOutput, output)
		}
	}
}

func TestMarshalCollection(t *testing.T) {
	loc, _ := time.LoadLocation("")
	schema := newMockSchema()

	tests := []struct {
		name          string
		data          Collection
		prepath       string
		params        string
		meta          map[string]interface{}
		jsonapi       map[string]interface{}
		errorExpected bool
		payloadFile   string
	}{
		{
			name: "collection with meta",
			data: mocktypes1,
			meta: map[string]interface{}{
				"num":       -32820,
				"timestamp": time.Date(1981, 2, 3, 4, 5, 6, 0, loc),
				"tf":        false,
				"str":       "//\n\téç.\\",
			},
			errorExpected: false,
			payloadFile:   "collection-1",
		}, {
			name:          "collection with prepath and params",
			data:          mocktypes2,
			prepath:       "https://example.org",
			params:        "?fields[mocktypes2]=uintptr,boolptr,timeptr",
			errorExpected: false,
			payloadFile:   "collection-2",
		}, {
			name:          "collection with prepath",
			data:          WrapCollection(Wrap(&mockType1{})),
			prepath:       "https://example.org",
			errorExpected: false,
			payloadFile:   "collection-3",
		},
	}

	for _, test := range tests {
		doc := NewDocument()
		doc.PrePath = test.prepath

		doc.Data = test.data

		resType := test.data.Type()
		rawurl := fmt.Sprintf("%s/%s%s", test.prepath, resType, test.params)

		url, err := ParseRawURL(schema, rawurl)
		tchek.UnintendedError(err)

		doc.Meta = test.meta

		// Marshal
		payload, err := Marshal(doc, url)
		tchek.ErrorExpected(t, test.name, test.errorExpected, err)

		if !test.errorExpected {
			var out bytes.Buffer

			// Format the payload
			json.Indent(&out, payload, "", "\t")
			output := out.String()

			// Retrieve the expected result from file
			content, err := ioutil.ReadFile("testdata/" + test.payloadFile + ".json")
			tchek.UnintendedError(err)
			out.Reset()
			json.Indent(&out, content, "", "\t")
			// Trim because otherwise there is an extra line at the end
			expectedOutput := strings.TrimSpace(out.String())

			tchek.AreEqual(t, test.name, expectedOutput, output)
		}
	}
}

func TestMarshalErrors(t *testing.T) {
	// Reset the IDs because the tests can't predict them.
	resetIDs := func(errors []Error) []Error {
		for i := range errors {
			errors[i].ID = "00000000-0000-0000-0000-000000000000"
		}
		return errors
	}

	tests := []struct {
		name          string
		errors        []Error
		errorExpected bool
		payloadFile   string
	}{
		{
			name: "two http errors",
			errors: resetIDs([]Error{
				NewErrBadRequest("Invalid attribute", "name cannot be empty."),
				NewErrBadRequest("Invalid attribute", "age cannot be negative."),
			}),
			errorExpected: false,
			payloadFile:   "errors-1",
		}, {
			name: "complex valid error",
			errors: resetIDs(func() []Error {
				e1 := NewError()

				e1.Code = "somecode"
				e1.Status = http.StatusInternalServerError
				e1.Title = "Error"
				e1.Detail = "An error occurred."
				e1.Links["about"] = "https://example.org/errors/about"
				e1.Source["pointer"] = "/data/attributes/title"
				e1.Meta["str"] = "a string"
				e1.Meta["num"] = 3943
				e1.Meta["bool"] = true

				return []Error{e1}
			}()),
			errorExpected: false,
			payloadFile:   "errors-2",
		},
	}

	for _, test := range tests {
		doc := NewDocument()
		doc.Data = test.errors
		// Marshal
		payload, err := Marshal(doc, nil)
		tchek.ErrorExpected(t, test.name, test.errorExpected, err)

		if !test.errorExpected {
			var out bytes.Buffer

			// Format the payload
			json.Indent(&out, payload, "", "\t")
			output := out.String()

			// Retrieve the expected result from file
			content, err := ioutil.ReadFile("testdata/" + test.payloadFile + ".json")
			tchek.UnintendedError(err)
			out.Reset()
			json.Indent(&out, content, "", "\t")
			// Trim because otherwise there is an extra line at the end
			expectedOutput := strings.TrimSpace(out.String())

			tchek.AreEqual(t, test.name, expectedOutput, output)
		}
	}
}
