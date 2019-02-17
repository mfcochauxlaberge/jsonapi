package jsonapi

import (
	"net/url"
	"testing"

	"github.com/mfcochauxlaberge/tchek"
)

func TestSimpleURL(t *testing.T) {
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
				PageSize:     10,
				PageNumber:   1,
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
				PageSize:     10,
				PageNumber:   1,
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
				PageSize:     10,
				PageNumber:   1,
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
				PageSize:     10,
				PageNumber:   1,
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
				PageSize:     10,
				PageNumber:   1,
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
				PageSize:     10,
				PageNumber:   1,
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
				PageSize:     10,
				PageNumber:   1,
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
				PageSize:     10,
				PageNumber:   1,
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
				PageSize:     10,
				PageNumber:   1,
				Include:      []string{},
			},
			expectedError: NewErrUnknownParameter("unknownparam"),
		},
	}

	for _, test := range tests {
		u, err := url.Parse(tchek.MakeOneLineNoSpaces(test.url))
		tchek.UnintendedError(err)

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

		tchek.AreEqual(t, test.name, test.expectedURL, url)
		tchek.AreEqual(t, test.name, test.expectedError, err)
	}
}

func TestParseCommaList(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedValue []string
	}{
		{
			name:          "empty",
			source:        ``,
			expectedValue: []string{},
		}, {
			name:          "comma only",
			source:        `,`,
			expectedValue: []string{},
		}, {
			name:          "two commas only",
			source:        `,,`,
			expectedValue: []string{},
		}, {
			name:          "single item",
			source:        `a`,
			expectedValue: []string{"a"},
		}, {
			name:          "start with comma",
			source:        `,a`,
			expectedValue: []string{"a"},
		}, {
			name:          "start with two commas",
			source:        `,,a`,
			expectedValue: []string{"a"},
		}, {
			name:          "start with comma and two items",
			source:        `,a,b`,
			expectedValue: []string{"a", "b"},
		}, {
			name:          "two items",
			source:        `a,b`,
			expectedValue: []string{"a", "b"},
		}, {
			name:          "two commas in middle",
			source:        `a,,b`,
			expectedValue: []string{"a", "b"},
		},
		{
			name:          "end with two commas",
			source:        `a,b,c,,`,
			expectedValue: []string{"a", "b", "c"},
		},
	}

	for _, test := range tests {
		value := parseCommaList(test.source)
		tchek.AreEqual(t, test.name, test.expectedValue, value)
	}
}

func TestParseFragments(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedValue []string
	}{
		{
			name:          "empty",
			source:        ``,
			expectedValue: []string{},
		}, {
			name:          "slash only",
			source:        `/`,
			expectedValue: []string{},
		}, {
			name:          "double slash",
			source:        `//`,
			expectedValue: []string{},
		}, {
			name:          "single item",
			source:        `a`,
			expectedValue: []string{"a"},
		}, {
			name:          "start with slash",
			source:        `/a`,
			expectedValue: []string{"a"},
		}, {
			name:          "start with two slashes",
			source:        `//a`,
			expectedValue: []string{"a"},
		}, {
			name:          "standard path",
			source:        `/a/b`,
			expectedValue: []string{"a", "b"},
		}, {
			name:          "two commas in middle",
			source:        `/a//b`,
			expectedValue: []string{"a", "b"},
		},
	}

	for _, test := range tests {
		value := parseFragments(test.source)
		tchek.AreEqual(t, test.name, test.expectedValue, value)
	}
}

func TestDeduceRoute(t *testing.T) {
	tests := []struct {
		name          string
		source        []string
		expectedValue string
	}{
		{
			name:          "empty",
			source:        []string{},
			expectedValue: "",
		}, {
			name:          "collection",
			source:        []string{"a"},
			expectedValue: "/a",
		}, {
			name:          "resource",
			source:        []string{"a", "b"},
			expectedValue: "/a/:id",
		}, {
			name:          "related relationship",
			source:        []string{"a", "b", "c"},
			expectedValue: "/a/:id/c",
		}, {
			name:          "self relationship",
			source:        []string{"a", "b", "relationships", "d"},
			expectedValue: "/a/:id/relationships/d",
		}, {
			name:          "collection meta",
			source:        []string{"a", "meta"},
			expectedValue: "/a/meta",
		}, {
			name:          "resource meta",
			source:        []string{"a", "b", "meta"},
			expectedValue: "/a/:id/meta",
		}, {
			name:          "related relationships meta",
			source:        []string{"a", "b", "relationships", "meta"},
			expectedValue: "/a/:id/relationships/meta",
		}, {
			name:          "self relationships meta",
			source:        []string{"a", "b", "relationships", "d", "meta"},
			expectedValue: "/a/:id/relationships/d/meta",
		},
	}

	for _, test := range tests {
		value := deduceRoute(test.source)
		tchek.AreEqual(t, test.name, test.expectedValue, value)
	}
}
