package jsonapi_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

// TestDocument ...
func TestDocument(t *testing.T) {
	assert := assert.New(t)

	pl1 := Document{}
	assert.Equal(nil, pl1.Data, "empty")
}

func TestInclude(t *testing.T) {
	assert := assert.New(t)

	doc := &Document{}
	typ1 := &Type{Name: "t1"}
	typ2 := &Type{Name: "t2"}

	/*
	 * Main data is a resource
	 */
	doc.Data = newResource(typ1, "id1")

	// Inclusions
	doc.Include(newResource(typ1, "id1"))
	doc.Include(newResource(typ1, "id2"))
	doc.Include(newResource(typ1, "id3"))
	doc.Include(newResource(typ1, "id3"))
	doc.Include(newResource(typ2, "id1"))

	// Check
	ids := []string{}
	for _, res := range doc.Included {
		ids = append(ids, res.GetType().Name+"-"+res.GetID())
	}

	expect := []string{
		"t1-id2",
		"t1-id3",
		"t2-id1",
	}
	assert.Equal(expect, ids)

	/*
	 * Main data is a collection
	 */
	doc = &Document{}

	// Collection
	col := &SoftCollection{}
	col.SetType(typ1)
	col.Add(newResource(typ1, "id1"))
	col.Add(newResource(typ1, "id2"))
	col.Add(newResource(typ1, "id3"))
	doc.Data = Collection(col)

	// Inclusions
	doc.Include(newResource(typ1, "id1"))
	doc.Include(newResource(typ1, "id2"))
	doc.Include(newResource(typ1, "id3"))
	doc.Include(newResource(typ1, "id4"))
	doc.Include(newResource(typ2, "id1"))
	doc.Include(newResource(typ2, "id1"))
	doc.Include(newResource(typ2, "id2"))

	// Check
	ids = []string{}
	for _, res := range doc.Included {
		ids = append(ids, res.GetType().Name+"-"+res.GetID())
	}

	expect = []string{
		"t1-id4",
		"t2-id1",
		"t2-id2",
	}
	assert.Equal(expect, ids)
}

func TestMarshalDocument(t *testing.T) {
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
	col.Add(Wrap(&mocktype{
		ID:    "id2",
		Str:   "漢語",
		Int:   -42,
		Time:  time.Time{},
		Bytes: []byte{},
	}))

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
			payload, err := MarshalDocument(test.doc, url)
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
				// TODO Figure out whether 0600 is okay or not.
				err = ioutil.WriteFile(path, dst.Bytes(), 0600)
				assert.NoError(err)
			}
		})
	}
}

func TestMarshalInvalidDocuments(t *testing.T) {
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
			_, err := MarshalDocument(test.doc, url)
			assert.EqualError(err, test.err)
		})
	}
}

func TestUnmarshalDocument(t *testing.T) {
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

		payload, err := MarshalDocument(doc, url)
		assert.NoError(err)

		doc2, err := UnmarshalDocument(payload, schema)
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

		payload, err := MarshalDocument(doc, url)
		assert.NoError(err)

		doc2, err := UnmarshalDocument(payload, schema)
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

		payload, err := MarshalDocument(doc, url)
		assert.NoError(err)

		doc2, err := UnmarshalDocument(payload, schema)
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
			doc, err := UnmarshalDocument([]byte(test.payload), schema)
			assert.EqualError(err, test.expected)
			assert.Nil(doc)
		}
	})
}

func newResource(typ *Type, id string) Resource {
	res := &SoftResource{}
	res.SetType(typ)
	res.SetID(id)

	return res
}
