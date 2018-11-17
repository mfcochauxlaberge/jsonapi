package jsonapi

import (
	"net/url"
	"testing"

	"github.com/kkaribu/tchek"
)

func TestSimpleURL(t *testing.T) {
	tests := []struct {
		url           string
		expectedURL   SimpleURL
		expectedError error
	}{

		{
			// 0
			url: ``,
			expectedURL: func() SimpleURL {
				sURL, _ := NewSimpleURL(nil)
				return sURL
			}(),
			expectedError: nil,
		}, {
			// 1
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
			// 2
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
			// 3
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
			// 4
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
			// 5
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
			// 6
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
					"type":  []string{"attr1", "attr2", "rel1"},
					"type2": []string{"attr3", "attr4", "rel2", "rel3"},
					"type3": []string{"attr5", "attr6", "rel4"},
					"type4": []string{"attr7", "rel5", "rel6"},
				},
				Filter:       nil,
				SortingRules: []string{"attr2", "-attr1"},
				PageSize:     20,
				PageNumber:   1,
				Include:      []string{"type2.type3", "type4"},
			},
			expectedError: nil,
		}, {
			// 7
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
			// 8
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
			// 9
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
			// 10
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

	for n, test := range tests {
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

		tchek.AreEqual(t, n, test.expectedURL, url)
		tchek.AreEqual(t, n, test.expectedError, err)
	}
}

func TestParseCommaList(t *testing.T) {
	tests := []struct {
		source        string
		expectedValue []string
	}{
		{
			// 0
			source:        ``,
			expectedValue: []string{},
		}, {
			// 1
			source:        `,`,
			expectedValue: []string{},
		}, {
			// 2
			source:        `,,`,
			expectedValue: []string{},
		}, {
			// 3
			source:        `a`,
			expectedValue: []string{"a"},
		}, {
			// 4
			source:        `,a`,
			expectedValue: []string{"a"},
		}, {
			// 5
			source:        `,,a`,
			expectedValue: []string{"a"},
		}, {
			// 6
			source:        `,a,b`,
			expectedValue: []string{"a", "b"},
		}, {
			// 7
			source:        `a,b`,
			expectedValue: []string{"a", "b"},
		}, {
			// 8
			source:        `a,,b`,
			expectedValue: []string{"a", "b"},
		},
		{
			// 8
			source:        `a,b,c,,`,
			expectedValue: []string{"a", "b", "c"},
		},
	}

	for n, test := range tests {
		value := parseCommaList(test.source)
		tchek.AreEqual(t, n, test.expectedValue, value)
	}
}

func TestParseFragments(t *testing.T) {
	tests := []struct {
		source        string
		expectedValue []string
	}{
		{
			// 0
			source:        ``,
			expectedValue: []string{},
		}, {
			// 1
			source:        `/`,
			expectedValue: []string{},
		}, {
			// 2
			source:        `//`,
			expectedValue: []string{},
		}, {
			// 3
			source:        `a`,
			expectedValue: []string{"a"},
		}, {
			// 4
			source:        `/a`,
			expectedValue: []string{"a"},
		}, {
			// 5
			source:        `//a`,
			expectedValue: []string{"a"},
		}, {
			// 6
			source:        `/a/b`,
			expectedValue: []string{"a", "b"},
		}, {
			// 7
			source:        `/a//b`,
			expectedValue: []string{"a", "b"},
		},
	}

	for n, test := range tests {
		value := parseFragments(test.source)
		tchek.AreEqual(t, n, test.expectedValue, value)
	}
}

func TestDeduceRoute(t *testing.T) {
	tests := []struct {
		source        []string
		expectedValue string
	}{
		{
			// 0
			source:        []string{},
			expectedValue: "",
		}, {
			// 1
			source:        []string{"a"},
			expectedValue: "/a",
		}, {
			// 2
			source:        []string{"a", "b"},
			expectedValue: "/a/:id",
		}, {
			// 3
			source:        []string{"a", "b", "c"},
			expectedValue: "/a/:id/c",
		}, {
			// 4
			source:        []string{"a", "b", "relationships", "d"},
			expectedValue: "/a/:id/relationships/d",
		}, {
			// 5
			source:        []string{"a", "meta"},
			expectedValue: "/a/meta",
		}, {
			// 6
			source:        []string{"a", "b", "meta"},
			expectedValue: "/a/:id/meta",
		}, {
			// 7
			source:        []string{"a", "b", "relationships", "meta"},
			expectedValue: "/a/:id/relationships/meta",
		}, {
			// 8
			source:        []string{"a", "b", "relationships", "d", "meta"},
			expectedValue: "/a/:id/relationships/d/meta",
		},
	}

	for n, test := range tests {
		value := deduceRoute(test.source)
		tchek.AreEqual(t, n, test.expectedValue, value)
	}
}
