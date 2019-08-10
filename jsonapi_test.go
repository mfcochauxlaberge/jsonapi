package jsonapi_test

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestMarshaling(t *testing.T) {
	// Update flag
	update := false
	flag.BoolVar(&update, "update-golden-files", false, "update the golden files")
	flag.Parse()

	// TODO Describe how this test suite works

	// Schema
	schema := getSchema()

	// Scenarios
	collections := []*SoftCollection{
		getmocktypeCollection(),
	}

	urls := []string{
		"/mocktype",
		"/mocktype/t1-1",
		"/mocktype/t1-1/relationships/to-1",
		"/mocktype/t1-1/relationships/to-x",
	}

	params := map[string][]string{
		"fields": []string{
			"?fields[mocktype]=id",
			"?fields[mocktype]=str",
			"?fields[mocktype]=to-1,to-x",
			"?fields[mocktype]=str,int,to-1,to-x",
		},
		"sort": []string{
			"&sort=str,int,id",
		},
		"pagination": []string{
			"",
			"&page[size]=0&page[number]=0",
			"&page[size]=2&page[number]=0",
			"&page[size]=1000&page[number]=0",
		},
	}

	lengths := []int{
		len(collections),
		len(urls),
		len(params["fields"]),
		len(params["sort"]),
		len(params["pagination"]),
	}

	// Test struct
	tests := []struct {
		name   string
		schema *Schema
		col    *SoftCollection
		url    string
	}{}

	counter := make([]int, len(lengths))
	run := true
	for run {
		col := collections[counter[0]]
		fullURL := urls[counter[1]] +
			params["fields"][counter[2]] +
			params["sort"][counter[3]] +
			params["pagination"][counter[4]]

		// Add test
		tests = append(tests, struct {
			name   string
			schema *Schema
			col    *SoftCollection
			url    string
		}{
			// TODO Give a different name to each test
			name:   "some name",
			schema: schema,
			col:    col,
			url:    fullURL,
		})

		// Increment counter
		for i := 0; i < len(counter); i++ {
			counter[i]++
			if counter[i] == lengths[i] {
				counter[i] = 0
				if i == len(counter)-1 {
					run = false
					break
				}
			} else {
				break
			}
		}
	}

	for i := range tests {
		i := i
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// URL
			url, err := NewURLFromRaw(test.schema, test.url)
			assert.NoError(err)
			fmt.Printf("URL: %s\n", url.UnescapedString())

			// Data
			var data interface{}
			if url.IsCol {
				// If it's a collection
				page := test.col.Range(nil, nil, nil, nil, 1000, 0)
				dataCol := &SoftCollection{
					Type: test.col.Type,
				}
				for i := range page {
					dataCol.Add(page[i])
				}
				data = dataCol
			} else {
				// If it's a resource
				for i := 0; i < test.col.Len(); i++ {
					if test.col.At(i).GetID() == url.ResID {
						data = test.col.At(i)
						break
					}
				}
			}

			// Document
			doc := &Document{
				Data: data,
			}

			// Marshaling
			payload, err := Marshal(doc, url)
			assert.NoError(err)

			// Golden file
			filename := "test" + strconv.Itoa(i) + ".json"
			path := filepath.Join("testdata", "goldenfiles", filename)
			if !update {
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

			// Unmarshaling
			// TODO
			// doc2, err := Unmarshal(payload, url, test.schema)
			// assert.NoError(err)
			// assert.Equal(doc, doc2)
		})
	}

	fmt.Printf("%d tests executed.\n", len(tests))
}

func getSchema() *Schema {
	schema := &Schema{}
	_ = schema.AddType(MustReflect(mocktype{}))
	if len(schema.Check()) > 0 {
		panic(" schema for tests should be valid")
	}
	return schema
}

func getEmptymocktypeCollection() *SoftCollection {
	schema := getSchema()
	typ := schema.GetType("mocktype")
	col := &SoftCollection{
		Type: &typ,
	}
	return col
}

func getmocktypeCollection() *SoftCollection {
	col := getEmptymocktypeCollection()
	col.Add(Wrap(&mocktype{
		ID: "t1-1",
	}))
	col.Add(Wrap(&mocktype{
		ID:       "t1-2",
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
		To1:      "",
		To1From1: "t1-10",
		To1FromX: "t1-11",
		ToX:      []string{},
		ToXFrom1: []string{"t1-12"},
		ToXFromX: []string{"t1-13", "t1-14"},
	}))
	return col
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

	// Relationships
	To1      string   `json:"to-1" api:"rel,mocktype"`
	To1From1 string   `json:"to-1-from-1" api:"rel,mocktype,to-1-from-1"`
	To1FromX string   `json:"to-1-from-x" api:"rel,mocktype,to-x-from-1"`
	ToX      []string `json:"to-x" api:"rel,mocktype"`
	ToXFrom1 []string `json:"to-x-from-1" api:"rel,mocktype,to-1-from-x"`
	ToXFromX []string `json:"to-x-from-x" api:"rel,mocktype,to-x-from-x"`
}
