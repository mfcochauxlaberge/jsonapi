package jsonapi_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

var update = flag.Bool("update-golden-files", false, "update the golden files")

func TestMarshaling(t *testing.T) {
	// TODO Describe how this test suite works

	// Setup
	typ, _ := BuildType(mocktype{})
	typ.NewFunc = func() Resource {
		return Wrap(&mocktype{})
	}
	col := &Resources{}
	col.Add(Wrap(&mocktype{
		ID:       "id1",
		Str:      "str",
		Int:      10,
		Int8:     18,
		Int16:    116,
		Int32:    132,
		Int64:    164,
		Uint:     100,
		Uint8:    108,
		Uint16:   1016,
		Uint32:   1032,
		Uint64:   1064,
		Bool:     true,
		Time:     getTime(),
		Bytes:    []byte{1, 2, 3},
		To1:      "id2",
		To1From1: "id3",
		To1FromX: "id3",
		ToX:      []string{"id2", "id3"},
		ToXFrom1: []string{"id4"},
		ToXFromX: []string{"id3", "id4"},
	}))
	col.Add(Wrap(&mocktype{
		ID:    "id2",
		Str:   "漢語",
		Int:   -42,
		Time:  time.Time{},
		Bytes: []byte{},
	}))
	col.Add(Wrap(&mocktype{ID: "id3"}))

	// Test struct
	tests := []struct {
		name   string
		doc    *Document
		fields []string
	}{
		{
			name: "empty data",
			doc: &Document{
				PrePath: "https://example.org",
			},
		}, {
			name: "empty collection",
			doc: &Document{
				Data: &Resources{},
			},
		}, {
			name: "resource",
			doc: &Document{
				Data: col.At(0),
				RelData: map[string][]string{
					"mocktype": {"to-1", "to-x-from-1"},
				},
			},
			fields: []string{
				"str", "uint64", "bool", "int", "time", "bytes", "to-1",
				"to-x-from-1",
			},
		}, {
			name: "collection",
			doc: &Document{
				Data: Range(col, nil, nil, []string{}, 10, 0),
				RelData: map[string][]string{
					"mocktype": {"to-1", "to-x-from-1"},
				},
				PrePath: "https://example.org",
			},
			fields: []string{
				"str", "uint64", "bool", "int", "time", "to-1", "to-x-from-1",
			},
		}, {
			name: "meta",
			doc: &Document{
				Data: nil,
				Meta: map[string]interface{}{
					"f1": "漢語",
					"f2": 42,
					"f3": true,
				},
			},
		}, {
			name: "collection with inclusions",
			doc: &Document{
				Data: Wrap(&mocktype{
					ID: "id1",
				}),
				RelData: map[string][]string{
					"mocktype": {"to-1", "to-x-from-1"},
				},
				Included: []Resource{
					Wrap(&mocktype{
						ID: "id2",
					}),
					Wrap(&mocktype{
						ID: "id3",
					}),
					Wrap(&mocktype{
						ID: "id4",
					}),
				},
			},
		}, {
			name: "identifier",
			doc: &Document{
				Data: Identifier{
					ID:   "id1",
					Type: "mocktype",
				},
			},
		}, {
			name: "identifiers",
			doc: &Document{
				Data: Identifiers{
					{
						ID:   "id1",
						Type: "mocktype",
					}, {
						ID:   "id2",
						Type: "mocktype",
					}, {
						ID:   "id3",
						Type: "mocktype",
					},
				},
			},
		}, {
			name: "error",
			doc: &Document{
				Errors: func() []Error {
					err := NewErrBadRequest("Bad Request", "This request is bad.")
					err.ID = "00000000-0000-0000-0000-000000000000"
					return []Error{err}
				}(),
			},
		}, {
			name: "errors",
			doc: &Document{
				Errors: func() []Error {
					err1 := NewErrBadRequest("Bad Request", "This request is bad.")
					err1.ID = "00000000-0000-0000-0000-000000000000"
					err2 := NewErrBadRequest("Bad Request", "This request is really bad.")
					err2.ID = "00000000-0000-0000-0000-000000000000"
					return []Error{err1, err2}
				}(),
			},
		},
	}

	for i := range tests {
		i := i
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// URL
			url := &URL{
				Fragments: []string{"fake", "path"},
				Params: &Params{
					Fields: map[string][]string{"mocktype": test.fields},
				},
			}
			if _, ok := test.doc.Data.(Collection); ok {
				url.IsCol = true
			}

			// Marshaling
			payload, err := Marshal(test.doc, url)
			assert.NoError(err)

			// Golden file
			filename := strings.Replace(test.name, " ", "_", -1) + ".json"
			path := filepath.Join("testdata", "goldenfiles", "marshaling", filename)
			if !*update {
				// Retrieve the expected result from file
				expected, _ := ioutil.ReadFile(path)
				assert.NoError(err, test.name)
				assert.JSONEq(string(expected), string(payload))
			} else {
				dst := &bytes.Buffer{}
				err = json.Indent(dst, payload, "", "\t")
				assert.NoError(err)
				// TODO Figure out whether 0644 is okay or not.
				err = ioutil.WriteFile(path, dst.Bytes(), 0644)
				assert.NoError(err)
			}
		})
	}
}

func TestMarshalingInvalidDocuments(t *testing.T) {
	// TODO Describe how this test suite works

	// Setup
	typ, _ := BuildType(mocktype{})
	typ.NewFunc = func() Resource {
		return Wrap(&mocktype{})
	}
	col := &Resources{}
	col.Add(Wrap(&mocktype{
		ID:       "id1",
		Str:      "str",
		Int:      10,
		Int8:     18,
		Int16:    116,
		Int32:    132,
		Int64:    164,
		Uint:     100,
		Uint8:    108,
		Uint16:   1016,
		Uint32:   1032,
		Uint64:   1064,
		Bool:     true,
		Time:     getTime(),
		To1:      "id2",
		To1From1: "id3",
		To1FromX: "id3",
		ToX:      []string{"id2", "id3"},
		ToXFrom1: []string{"id4"},
		ToXFromX: []string{"id3", "id4"},
	}))
	col.Add(Wrap(&mocktype{
		ID:   "id2",
		Str:  "漢語",
		Int:  -42,
		Time: time.Time{},
	}))
	col.Add(Wrap(&mocktype{ID: "id3"}))

	// Test struct
	tests := []struct {
		name   string
		doc    *Document
		fields []string
		err    string
	}{
		{
			name: "invalid data",
			doc: &Document{
				Data: "just a string",
			},
			err: "data contains an unknown type",
		},
	}

	for i := range tests {
		i := i
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// URL
			url := &URL{
				Fragments: []string{"fake", "path"},
				Params: &Params{
					Fields: map[string][]string{"mocktype": test.fields},
				},
			}
			if _, ok := test.doc.Data.(Collection); ok {
				url.IsCol = true
			}

			// Marshaling
			_, err := Marshal(test.doc, url)
			assert.EqualError(err, test.err)
		})
	}
}

func TestUnmarshaling(t *testing.T) {
	// Setup
	typ, _ := BuildType(mocktype{})
	typ.NewFunc = func() Resource {
		return Wrap(&mocktype{})
	}
	schema := &Schema{Types: []Type{typ}}
	col := Resources{}
	col.Add(Wrap(&mocktype{
		ID:       "id1",
		Str:      "str",
		Int:      10,
		Int8:     18,
		Int16:    116,
		Int32:    132,
		Int64:    164,
		Uint:     100,
		Uint8:    108,
		Uint16:   1016,
		Uint32:   1032,
		Uint64:   1064,
		Bool:     true,
		Time:     getTime(),
		Bytes:    []byte{1, 2, 3},
		To1:      "id2",
		To1From1: "id3",
		To1FromX: "id3",
		ToX:      []string{"id2", "id3"},
		ToXFrom1: []string{"id4"},
		ToXFromX: []string{"id3", "id4"},
	}))
	col.Add(Wrap(&mocktype{ID: "id2"}))
	col.Add(Wrap(&mocktype{ID: "id3"}))

	// Tests
	t.Run("resource with inclusions", func(t *testing.T) {
		assert := assert.New(t)

		url, _ := NewURLFromRaw(schema, "/mocktype/id1")

		doc := &Document{
			Data: col.At(0),
			RelData: map[string][]string{
				"mocktype": typ.Fields(),
			},
			Included: []Resource{
				col.At(1),
				col.At(2),
			},
		}

		payload, err := Marshal(doc, url)
		assert.NoError(err)

		doc2, err := Unmarshal(payload, schema)
		assert.NoError(err)
		assert.True(Equal(doc.Data.(Resource), doc2.Data.(Resource)))
		// TODO Make all the necessary assertions.
	})

	t.Run("collection with inclusions", func(t *testing.T) {
		assert := assert.New(t)

		url, _ := NewURLFromRaw(schema, "/mocktype/id1")

		doc := &Document{
			Data: &col,
			RelData: map[string][]string{
				"mocktype": typ.Fields(),
			},
		}

		payload, err := Marshal(doc, url)
		assert.NoError(err)

		doc2, err := Unmarshal(payload, schema)
		assert.NoError(err)
		assert.IsType(&col, doc.Data)
		assert.IsType(&col, doc2.Data)
		if col, ok := doc.Data.(Collection); ok {
			if col2, ok := doc2.Data.(Collection); ok {
				assert.Equal(col.Len(), col2.Len())
				for j := 0; j < col.Len(); j++ {
					assert.True(Equal(col.At(j), col2.At(j)))
				}
			}
		}
		// TODO Make all the necessary assertions.
	})

	t.Run("identifier", func(t *testing.T) {
		assert := assert.New(t)

		url, _ := NewURLFromRaw(schema, "/mocktype/id1/relationships/to-1")

		doc := &Document{
			Data: Identifier{
				ID:   "id2",
				Type: "mocktype",
			},
		}

		payload, err := Marshal(doc, url)
		assert.NoError(err)

		doc2, err := UnmarshalIdentifiers(payload, schema)
		assert.NoError(err)
		assert.Equal(doc.Data, doc2.Data)
	})

	t.Run("identifers", func(t *testing.T) {
		assert := assert.New(t)

		url, _ := NewURLFromRaw(schema, "/mocktype/id1/relationships/to-x")

		doc := &Document{
			Data: Identifiers{
				Identifier{
					ID:   "id2",
					Type: "mocktype",
				},
				Identifier{
					ID:   "id3",
					Type: "mocktype",
				},
			},
		}

		payload, err := Marshal(doc, url)
		assert.NoError(err)

		doc2, err := UnmarshalIdentifiers(payload, schema)
		assert.NoError(err)
		assert.Equal(doc.Data, doc2.Data)
	})

	t.Run("errors (Unmarshal)", func(t *testing.T) {
		assert := assert.New(t)

		url, _ := NewURLFromRaw(schema, "/mocktype/id1/relationships/to-x")

		doc := &Document{
			Errors: func() []Error {
				err := NewErrBadRequest("Bad Request", "This request is bad.")
				err.ID = "00000000-0000-0000-0000-000000000000"
				return []Error{err}
			}(),
		}

		payload, err := Marshal(doc, url)
		assert.NoError(err)

		doc2, err := Unmarshal(payload, schema)
		assert.NoError(err)
		assert.Equal(doc.Data, doc2.Data)
	})

	t.Run("errors (UnmarshalIdentifers)", func(t *testing.T) {
		assert := assert.New(t)

		url, _ := NewURLFromRaw(schema, "/mocktype/id1/relationships/to-x")

		doc := &Document{
			Errors: func() []Error {
				err := NewErrBadRequest("Bad Request", "This request is bad.")
				err.ID = "00000000-0000-0000-0000-000000000000"
				return []Error{err}
			}(),
		}

		payload, err := Marshal(doc, url)
		assert.NoError(err)

		doc2, err := UnmarshalIdentifiers(payload, schema)
		assert.NoError(err)
		assert.Equal(doc.Data, doc2.Data)
	})

	t.Run("invalid payloads (Unmarshal)", func(t *testing.T) {
		assert := assert.New(t)

		tests := []struct {
			payload  string
			expected string
		}{
			{
				payload:  `invalid payload`,
				expected: "invalid character 'i' looking for beginning of value",
			}, {
				payload:  `{"data":"invaliddata"}`,
				expected: "400 Bad Request: Missing data top-level member in payload.",
			}, {
				payload:  `{"data":{"id":true}}`,
				expected: "400 Bad Request: The provided JSON body could not be read.",
			}, {
				payload:  `{"data":[{"id":true}]}`,
				expected: "400 Bad Request: The provided JSON body could not be read.",
			}, {
				payload:  `{"jsonapi":{"key":"data/errors missing"}}`,
				expected: "400 Bad Request: Missing data top-level member in payload.",
			}, {
				payload: `{"data":null,"included":[{"id":true}]}`,
				expected: "json: " +
					"cannot unmarshal bool into Go struct field Identifier.id of type string",
			}, {
				payload:  `{"data":null,"included":[{"attributes":true}]}`,
				expected: "400 Bad Request: The provided JSON body could not be read.",
			}, {
				payload:  `{"data":{"id":"1","type":"mocktype","attributes":{"nonexistent":1}}}`,
				expected: "400 Bad Request: \"nonexistent\" is not a known field.",
			}, {
				payload:  `{"data":{"id":"1","type":"mocktype","attributes":{"int8":"abc"}}}`,
				expected: "400 Bad Request: The field value is invalid for the expected type.",
			}, {
				payload: `{
					"data": {
						"id": "1",
						"type": "mocktype",
						"relationships": {
							"to-x": {
								"data": "wrong"
							}
						}
					}
				}`,
				expected: "400 Bad Request: The field value is invalid for the expected type.",
			}, {
				payload: `{
					"data": {
						"id": "1",
						"type": "mocktype",
						"relationships": {
							"wrong": {
								"data": "wrong"
							}
						}
					}
				}`,
				expected: "400 Bad Request: \"wrong\" is not a known field.",
			},
		}

		for _, test := range tests {
			doc, err := Unmarshal([]byte(test.payload), schema)
			assert.EqualError(err, test.expected)
			assert.Nil(doc)
		}
	})

	t.Run("invalid payloads (UnmarshalIdentifiers)", func(t *testing.T) {
		assert := assert.New(t)

		tests := []struct {
			payload  string
			expected string
		}{
			{
				payload:  `{invalid json}`,
				expected: "invalid character 'i' looking for beginning of object key string",
			}, {
				payload:  `{"jsonapi":{}}`,
				expected: "400 Bad Request: Missing data top-level member in payload.",
			}, {
				payload:  `{"jsonapi":{"key":"data/errors missing"}}`,
				expected: "400 Bad Request: Missing data top-level member in payload.",
			}, {
				payload: `{"data":{"id":["invalid"]}}`,
				expected: "json: " +
					"cannot unmarshal array into Go struct field Identifier.id of type string",
			}, {
				payload: `{"data":[{"id":["invalid"]}]}`,
				expected: "json: " +
					"cannot unmarshal array into Go struct field Identifier.id of type string",
			},
		}

		for _, test := range tests {
			doc, err := UnmarshalIdentifiers([]byte(test.payload), nil)
			assert.EqualError(err, test.expected)
			assert.Nil(doc)
		}
	})
}

func getTime() time.Time {
	now, _ := time.Parse(time.RFC3339Nano, "2013-06-24T22:03:34.8276Z")
	return now
}

// mocktype is a fake struct that defines a JSON:API type for test purposes.
type mocktype struct {
	ID string `json:"id" api:"mocktype"`

	// Attributes
	Str    string    `json:"str" api:"attr"`
	Int    int       `json:"int" api:"attr"`
	Int8   int8      `json:"int8" api:"attr"`
	Int16  int16     `json:"int16" api:"attr"`
	Int32  int32     `json:"int32" api:"attr"`
	Int64  int64     `json:"int64" api:"attr"`
	Uint   uint      `json:"uint" api:"attr"`
	Uint8  uint8     `json:"uint8" api:"attr"`
	Uint16 uint16    `json:"uint16" api:"attr"`
	Uint32 uint32    `json:"uint32" api:"attr"`
	Uint64 uint64    `json:"uint64" api:"attr"`
	Bool   bool      `json:"bool" api:"attr"`
	Time   time.Time `json:"time" api:"attr"`
	Bytes  []byte    `json:"bytes" api:"attr"`

	// Relationships
	To1      string   `json:"to-1" api:"rel,mocktype"`
	To1From1 string   `json:"to-1-from-1" api:"rel,mocktype,to-1-from-1"`
	To1FromX string   `json:"to-1-from-x" api:"rel,mocktype,to-x-from-1"`
	ToX      []string `json:"to-x" api:"rel,mocktype"`
	ToXFrom1 []string `json:"to-x-from-1" api:"rel,mocktype,to-1-from-x"`
	ToXFromX []string `json:"to-x-from-x" api:"rel,mocktype,to-x-from-x"`
}
