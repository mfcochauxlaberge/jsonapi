package jsonapi

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"testing"
// 	"time"

// 	. "github.com/mfcochauxlaberge/jsonapi"

// 	"github.com/stretchr/testify/assert"
// )

// func TestMarshalResource(t *testing.T) {
// 	assert := assert.New(t)

// 	loc, _ := time.LoadLocation("")
// 	schema := newMockSchema()

// 	tests := []struct {
// 		name          string
// 		data          Resource
// 		inclusions    []Resource
// 		relData       map[string][]string
// 		prepath       string
// 		params        string
// 		meta          map[string]interface{}
// 		errorExpected bool
// 		payloadFile   string
// 	}{
// 		{
// 			name: "resource with meta",
// 			data: mocktypes1.At(0),
// 			meta: map[string]interface{}{
// 				"num":       42,
// 				"timestamp": time.Date(2017, 1, 2, 3, 4, 5, 6, loc),
// 				"tf":        true,
// 				"str":       "a string",
// 			},
// 			errorExpected: false,
// 			payloadFile:   "resource-1",
// 		}, {
// 			name:          "resource with prepath",
// 			data:          mocktypes2.At(1),
// 			prepath:       "https://example.org",
// 			errorExpected: false,
// 			payloadFile:   "resource-2",
// 		}, {
// 			name:          "resource with prepath and params",
// 			data:          mocktypes2.At(1),
// 			prepath:       "https://example.org",
// 			params:        "?fields[mocktypes2]=strptr,uintptr,int",
// 			errorExpected: false,
// 			payloadFile:   "resource-3",
// 		}, {
// 			name:          "resource with no attributes and relationships",
// 			data:          mocktypes1.At(0),
// 			prepath:       "https://example.org",
// 			params:        "?fields[mocktypes1]=id",
// 			errorExpected: false,
// 			payloadFile:   "resource-5",
// 		}, {
// 			name: "resource with relationship data",
// 			data: mocktypes11.At(0),
// 			relData: map[string][]string{
// 				"mocktypes1": []string{
// 					"to-one", "to-many", "to-one-from-one", "to-many-from-many",
// 				},
// 			},
// 			prepath: "https://example.org",
// 			params: `
// 				?fields[mocktypes1]=
// 					to-one,to-many,
// 					to-one-from-one,to-one-from-many,
// 					to-many-from-one,to-many-from-many
// 			`,
// 			errorExpected: false,
// 			payloadFile:   "resource-6",
// 		}, {
// 			name: "resource with inclusions",
// 			data: mocktypes11.At(0),
// 			inclusions: []Resource{
// 				mocktypes21.At(0),
// 				mocktypes21.At(1),
// 				mocktypes21.At(2),
// 			},
// 			relData: map[string][]string{
// 				"mocktypes1": []string{
// 					"to-one", "to-many", "to-one-from-one", "to-many-from-many",
// 				},
// 			},
// 			prepath: "https://example.org",
// 			params: `
// 				?fields[mocktypes1]=
// 					to-one,to-many
// 				&fields[mocktypes2]=
// 					intptr,boolptr,strptr,
// 			`,
// 			errorExpected: false,
// 			payloadFile:   "resource-7",
// 		},
// 	}

// 	for _, test := range tests {
// 		doc := NewDocument()
// 		doc.PrePath = test.prepath

// 		doc.Data = test.data

// 		id := test.data.GetID()
// 		typ := test.data.GetType()
// 		rawurl := fmt.Sprintf(
// 			"%s/%s/%s%s",
// 			test.prepath, typ.Name, id, makeOneLineNoSpaces(test.params),
// 		)

// 		url, err := NewURLFromRaw(schema, rawurl)
// 		assert.NoError(err, test.name)

// 		for _, inc := range test.inclusions {
// 			doc.Include(inc)
// 		}

// 		doc.RelData = test.relData
// 		doc.Meta = test.meta

// 		// Marshal
// 		payload, err := Marshal(doc, url)

// 		if test.errorExpected {
// 			assert.Error(err, test.name)
// 		} else {
// 			assert.NoError(err, test.name)
// 			// Retrieve the expected result from file
// 			expected, _ := ioutil.ReadFile("testdata/" + test.payloadFile + ".json")
// 			assert.NoError(err, test.name)
// 			assert.JSONEq(string(expected), string(payload), test.name)
// 		}
// 	}
// }

// func TestMarshalCollection(t *testing.T) {
// 	assert := assert.New(t)

// 	loc, _ := time.LoadLocation("")
// 	schema := newMockSchema()

// 	tests := []struct {
// 		name          string
// 		data          Collection
// 		prepath       string
// 		params        string
// 		meta          map[string]interface{}
// 		errorExpected bool
// 		payloadFile   string
// 	}{
// 		{
// 			name: "collection with meta",
// 			data: mocktypes1,
// 			meta: map[string]interface{}{
// 				"num":       -32820,
// 				"timestamp": time.Date(1981, 2, 3, 4, 5, 6, 0, loc),
// 				"tf":        false,
// 				"str":       "//\n\téç.\\",
// 			},
// 			errorExpected: false,
// 			payloadFile:   "collection-1",
// 		}, {
// 			name:          "collection with prepath and params",
// 			data:          mocktypes2,
// 			prepath:       "https://example.org",
// 			params:        "?fields[mocktypes2]=uintptr,boolptr,timeptr",
// 			errorExpected: false,
// 			payloadFile:   "collection-2",
// 		}, {
// 			name:          "collection with prepath",
// 			data:          WrapCollection(Wrap(&mockType1{})),
// 			prepath:       "https://example.org",
// 			errorExpected: false,
// 			payloadFile:   "collection-3",
// 		},
// 	}

// 	for _, test := range tests {
// 		doc := NewDocument()
// 		doc.PrePath = test.prepath

// 		doc.Data = test.data

// 		typ := test.data.GetType()
// 		rawurl := fmt.Sprintf("%s/%s%s", test.prepath, typ.Name, test.params)

// 		url, err := NewURLFromRaw(schema, rawurl)
// 		assert.NoError(err, test.name)

// 		doc.Meta = test.meta

// 		// Marshal
// 		payload, err := Marshal(doc, url)

// 		if test.errorExpected {
// 			assert.Error(err, test.name)
// 		} else {
// 			assert.NoError(err, test.name)
// 			// Retrieve the expected result from file
// 			expected, _ := ioutil.ReadFile("testdata/" + test.payloadFile + ".json")
// 			assert.JSONEq(string(expected), string(payload), test.name)
// 		}
// 	}
// }

// func TestMarshalInclusions(t *testing.T) {
// 	assert := assert.New(t)

// 	schema := newMockSchema()

// 	// Document
// 	doc := &Document{}
// 	doc.PrePath = "https://example.org"

// 	// URL
// 	url, err := NewURLFromRaw(
// 		schema,
// 		makeOneLineNoSpaces(`
// 			/mocktypes3/mt3-1
// 			?fields[mocktypes1]=str
// 			&fields[mocktypes3]=attr1,attr2
// 		`),
// 	)
// 	assert.NoError(err)

// 	// Data (single resource)
// 	res := Wrap(&mockType3{})
// 	res.SetID("mt3-1")
// 	res.Set("attr1", "str")
// 	res.Set("attr2", 42)
// 	doc.Data = Resource(res)

// 	// Inclusions
// 	inc1 := Wrap(&mockType1{})
// 	inc1.SetID("mt1-1")
// 	inc1.Set("str", "astring")
// 	doc.Include(inc1)

// 	inc2 := Wrap(&mockType1{})
// 	inc2.SetID("mt1-2")
// 	inc2.Set("str", "anotherstring")
// 	doc.Include(inc2)

// 	payload, _ := Marshal(doc, url)

// 	// Retrieve the expected result from file
// 	expected, _ := ioutil.ReadFile("testdata/resource-4.json")

// 	assert.JSONEq(string(expected), string(payload))
// }

// func TestMarshalErrors(t *testing.T) {
// 	assert := assert.New(t)

// 	// Reset the IDs because the tests can't predict them.
// 	resetIDs := func(errors []Error) []Error {
// 		for i := range errors {
// 			errors[i].ID = "00000000-0000-0000-0000-000000000000"
// 		}
// 		return errors
// 	}

// 	tests := []struct {
// 		name          string
// 		errors        []Error
// 		errorExpected bool
// 		payloadFile   string
// 	}{
// 		{
// 			name: "two http errors",
// 			errors: resetIDs([]Error{
// 				NewErrBadRequest("Invalid attribute", "name cannot be empty."),
// 				NewErrBadRequest("Invalid attribute", "age cannot be negative."),
// 			}),
// 			errorExpected: false,
// 			payloadFile:   "errors-1",
// 		}, {
// 			name: "complex valid error",
// 			errors: resetIDs(func() []Error {
// 				e1 := NewError()

// 				e1.Code = "somecode"
// 				e1.Status = http.StatusInternalServerError
// 				e1.Title = "Error"
// 				e1.Detail = "An error occurred."
// 				e1.Links["about"] = "https://example.org/errors/about"
// 				e1.Source["pointer"] = "/data/attributes/title"
// 				e1.Meta["str"] = "a string"
// 				e1.Meta["num"] = 3943
// 				e1.Meta["bool"] = true

// 				return []Error{e1}
// 			}()),
// 			errorExpected: false,
// 			payloadFile:   "errors-2",
// 		},
// 	}

// 	for _, test := range tests {
// 		doc := NewDocument()
// 		if len(test.errors) == 1 {
// 			doc.Data = test.errors[0]
// 		} else {
// 			doc.Data = test.errors
// 		}

// 		// Marshal
// 		payload, err := Marshal(doc, nil)

// 		if test.errorExpected {
// 			assert.Error(err, test.name)
// 		} else {
// 			assert.NoError(err, test.name)
// 			// Retrieve the expected result from file
// 			expected, _ := ioutil.ReadFile("testdata/" + test.payloadFile + ".json")
// 			assert.JSONEq(string(expected), string(payload), test.name)
// 		}
// 	}
// }

// func TestMarshalOther(t *testing.T) {
// 	assert := assert.New(t)

// 	doc := &Document{
// 		Data: nil,
// 	}
// 	payload, err := Marshal(doc, nil)
// 	assert.NoError(err)

// 	// Retrieve the expected result from file
// 	expected, _ := ioutil.ReadFile("testdata/null-1.json")

// 	assert.JSONEq(string(expected), string(payload), "null data")
// }
