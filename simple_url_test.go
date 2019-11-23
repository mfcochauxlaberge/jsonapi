package jsonapi_test

import (
	"net/url"
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestSimpleURL(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name          string
		url           string
		expectedURL   SimpleURL
		expectedError error
	}{

		{
			name: "empty url",
			url:  ``,
			expectedURL: func() SimpleURL {
				sURL, _ := NewSimpleURL(nil)
				return sURL
			}(),
			expectedError: nil,
		}, {
			name: "empty path",
			url: `
				http://api.example.com
			`,
			expectedURL: SimpleURL{
				Fragments: []string{},
				Route:     "",

				Fields:       map[string][]string{},
				Filter:       nil,
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: nil,
		}, {
			name: "collection",
			url: `
				http://api.example.com/type
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type"},
				Route:     "/type",

				Fields:       map[string][]string{},
				Filter:       nil,
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: nil,
		}, {
			name: "resource",
			url: `
				http://api.example.com/type/id
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type", "id"},
				Route:     "/type/:id",

				Fields:       map[string][]string{},
				Filter:       nil,
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: nil,
		}, {
			name: "relationship",
			url: `
				http://api.example.com/type/id/rel
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type", "id", "rel"},
				Route:     "/type/:id/rel",

				Fields:       map[string][]string{},
				Filter:       nil,
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: nil,
		}, {
			name: "self relationship",
			url: `
				http://api.example.com/type/id/relationships/rel
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type", "id", "relationships", "rel"},
				Route:     "/type/:id/relationships/rel",

				Fields:       map[string][]string{},
				Filter:       nil,
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: nil,
		}, {
			name: "fields, sort, pagination, include",
			url: `
				http://api.example.com/type
				?fields[type]=attr1,attr2,rel1
				&fields[type2]=attr3,attr4,rel2,rel3
				&fields[type3]=attr5,attr6,rel4
				&fields[type4]=attr7,rel5,rel6
				&sort=attr2,-attr1
				&page[number]=1
				&page[size]=20
				&include=type2.type3,type4
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type"},
				Route:     "/type",

				Fields: map[string][]string{
					"type":  {"attr1", "attr2", "rel1"},
					"type2": {"attr3", "attr4", "rel2", "rel3"},
					"type3": {"attr5", "attr6", "rel4"},
					"type4": {"attr7", "rel5", "rel6"},
				},
				Filter:       nil,
				SortingRules: []string{"attr2", "-attr1"},
				PageSize:     20,
				PageNumber:   1,
				Include:      []string{"type2.type3", "type4"},
			},
			expectedError: nil,
		}, {
			name: "filter label",
			url: `
				http://api.example.com/type/id/rel
				?filter=label
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type", "id", "rel"},
				Route:     "/type/:id/rel",

				Fields:       map[string][]string{},
				FilterLabel:  "label",
				Filter:       nil,
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: nil,
		}, {
			name: "negative page size",
			url: `
				http://api.example.com/type/id/rel
				?page[size]=-1
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type", "id", "rel"},
				Route:     "/type/:id/rel",

				Fields:       map[string][]string{},
				Filter:       nil,
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: NewErrInvalidPageSizeParameter("-1"),
		}, {
			name: "negative page number",
			url: `
				http://api.example.com/type/id/rel
				?page[number]=-1
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type", "id", "rel"},
				Route:     "/type/:id/rel",

				Fields:       map[string][]string{},
				Filter:       nil,
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: NewErrInvalidPageNumberParameter("-1"),
		}, {
			name: "unknown parameter",
			url: `
				http://api.example.com/type/id/rel
				?unknownparam=somevalue
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type", "id", "rel"},
				Route:     "/type/:id/rel",

				Fields:       map[string][]string{},
				Filter:       nil,
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: NewErrUnknownParameter("unknownparam"),
		}, {
			name: "filter query",
			url: `
				http://api.example.com/type/id/rel
				?filter={
					"f": "field",
					"o": "=",
					"v": "abc"
				}
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type", "id", "rel"},
				Route:     "/type/:id/rel",

				Fields: map[string][]string{},
				Filter: &Filter{
					Field: "field",
					Op:    "=",
					Val:   "abc",
				},
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: nil,
		}, {
			name: "filter query",
			url: `
				http://api.example.com/type/id/rel
				?filter={"thisis:invalid"}
			`,
			expectedURL: SimpleURL{
				Fragments: []string{"type", "id", "rel"},
				Route:     "/type/:id/rel",

				Fields:       map[string][]string{},
				Filter:       nil,
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      []string{},
			},
			expectedError: NewErrMalformedFilterParameter(`{"thisis:invalid"}`),
		},
	}

	for _, test := range tests {
		u, err := url.Parse(makeOneLineNoSpaces(test.url))
		assert.NoError(err, test.name)

		url, err := NewSimpleURL(u)

		if test.expectedError != nil {
			jaErr := test.expectedError.(Error)
			jaErr.ID = ""
			test.expectedError = jaErr
		}

		if err != nil {
			jaErr := err.(Error)
			jaErr.ID = ""
			err = jaErr
		}

		assert.Equal(test.expectedURL, url, test.name)
		assert.Equal(test.expectedError, err, test.name)
	}
}

func TestSimpleURLPath(t *testing.T) {
	assert := assert.New(t)

	su := &SimpleURL{Fragments: []string{}}
	assert.Equal("", su.Path())

	su = &SimpleURL{Fragments: []string{"a", "b", "c"}}
	assert.Equal("a/b/c", su.Path())
}
