package jsonapi

import (
	"net/url"
	"testing"

	"github.com/kkaribu/tchek"
)

func TestParseURL(t *testing.T) {
	// Registry
	reg := NewMockRegistry()

	tests := []struct {
		url           string
		expectedURL   URL
		expectedError bool
	}{

		{
			// 0
			url:           ``,
			expectedURL:   URL{},
			expectedError: true,
		}, {
			// 1
			url:           `https://example.com`,
			expectedURL:   URL{},
			expectedError: true,
		}, {
			// 2
			url: `mocktypes1`,
			expectedURL: URL{
				URL:           "mocktypes1",
				URLNormalized: "/mocktypes1",
				Path:          []string{"mocktypes1"},
				Route:         "/mocktypes1",
				FromFilter:    FromFilter{},
				ResType:       "mocktypes1",
				ResID:         "",
				RelKind:       "",
			},
			expectedError: false,
		}, {
			// 3
			url: `https://example.com/mocktypes1`,
			expectedURL: URL{
				URL:           "/mocktypes1",
				URLNormalized: "/mocktypes1",
				Path:          []string{"mocktypes1"},
				Route:         "/mocktypes1",
				FromFilter:    FromFilter{},
				ResType:       "mocktypes1",
				ResID:         "",
				RelKind:       "",
			},
			expectedError: false,
		}, {
			// 4
			url: `https://example.com/mocktypes1/mc1-1`,
			expectedURL: URL{
				URL:           "/mocktypes1/mc1-1",
				URLNormalized: "/mocktypes1/mc1-1",
				Path:          []string{"mocktypes1", "mc1-1"},
				Route:         "/mocktypes1/:id",
				FromFilter:    FromFilter{},
				ResType:       "mocktypes1",
				ResID:         "mc1-1",
				RelKind:       "",
			},
			expectedError: false,
		}, {
			// 5
			url: `https://example.com/mocktypes1/mc1-1/to-one`,
			expectedURL: URL{
				URL:           "/mocktypes1/mc1-1/to-one",
				URLNormalized: "/mocktypes1/mc1-1/to-one",
				Path:          []string{"mocktypes1", "mc1-1", "to-one"},
				Route:         "/mocktypes1/:id/to-one",
				FromFilter: FromFilter{
					Type: "mocktypes1",
					ID:   "mc1-1",
					Name: "to-one",
				},
				ResType: "mocktypes2",
				ResID:   "",
				RelKind: "related",
			},
			expectedError: false,
		}, {
			// 6
			url: `https://example.com/mocktypes1/mc1-1/relationships/to-many-from-one`,
			expectedURL: URL{
				URL:           "/mocktypes1/mc1-1/relationships/to-many-from-one",
				URLNormalized: "/mocktypes1/mc1-1/relationships/to-many-from-one",
				Path:          []string{"mocktypes1", "mc1-1", "relationships", "to-many-from-one"},
				Route:         "/mocktypes1/:id/relationships/to-many-from-one",
				FromFilter: FromFilter{
					Type:        "mocktypes1",
					ID:          "mc1-1",
					Name:        "to-many-from-one",
					InverseName: "to-one-from-many",
				},
				ResType: "mocktypes2",
				ResID:   "",
				RelKind: "self",
			},
			expectedError: false,
		}, {
			// 7
			url: `/mocktypes1/mc1-1/relationships/to-many-from-one`,
			expectedURL: URL{
				URL:           "/mocktypes1/mc1-1/relationships/to-many-from-one",
				URLNormalized: "/mocktypes1/mc1-1/relationships/to-many-from-one",
				Path:          []string{"mocktypes1", "mc1-1", "relationships", "to-many-from-one"},
				Route:         "/mocktypes1/:id/relationships/to-many-from-one",
				FromFilter: FromFilter{
					Type:        "mocktypes1",
					ID:          "mc1-1",
					Name:        "to-many-from-one",
					InverseName: "to-one-from-many",
				},
				ResType: "mocktypes2",
				ResID:   "",
				RelKind: "self",
			},
			expectedError: false,
		},
	}

	for n, test := range tests {
		u, _ := url.Parse(tchek.MakeOneLineNoSpaces(test.url))
		url, err := ParseURL(reg, u)
		tchek.ErrorExpected(t, n, test.expectedError, err)

		if !test.expectedError {
			url.Params = nil
			tchek.AreEqual(
				t, n,
				test.expectedURL,
				*url,
			)
		}
	}
}

func TestParseParams(t *testing.T) {
	// Registry
	reg := NewMockRegistry()

	tests := []struct {
		url            string
		resType        string
		expectedParams Params
		expectedError  bool
	}{

		{
			// 0
			url: `/`,
			expectedParams: Params{
				Fields:       map[string][]string{"": nil},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				AttrFilters:  map[string]AttrFilter{},
				RelFilters:   map[string]RelFilter{},
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      [][]Rel{},
			},
			expectedError: false,
		}, {
			// 1
			url: `?`,
			expectedParams: Params{
				Fields:       map[string][]string{"": nil},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				AttrFilters:  map[string]AttrFilter{},
				RelFilters:   map[string]RelFilter{},
				SortingRules: []string{},
				PageSize:     0,
				PageNumber:   0,
				Include:      [][]Rel{},
			},
			expectedError: false,
		},
		// 	// 2
		// 	url: `
		// 		?include=
		// 			to-many-from-one.to-one-from-many.to-one.to-many,
		// 			to-one-from-one.to-many-from-many
		// 		&include=
		// 			to-many-from-one,
		// 			to-many-from-many
		// 	`,
		// 	resType: "mocktypes1",
		// 	expectedParams: Params{
		// 		Fields:       map[string][]string{},
		// 		Attrs:        map[string][]Attr{},
		// 		Rels:         map[string][]Rel{},
		// 		RelData:      map[string][]string{},
		// 		AttrFilters:  map[string]AttrFilter{},
		// 		RelFilters:   map[string]RelFilter{},
		// 		SortingRules: []string{},
		// 		PageSize:     0,
		// 		PageNumber:   0,
		// 		Include:      [][]Rel{},
		// 	},
		// 	expectedError: false,
		// },
	}

	for n, test := range tests {
		u, err := url.Parse(tchek.MakeOneLineNoSpaces(test.url))
		tchek.UnintendedError(err)

		params, err := parseParams(reg, test.resType, u)
		tchek.ErrorExpected(t, n, test.expectedError, err)

		if !test.expectedError {
			tchek.AreEqual(
				t, n,
				test.expectedParams,
				*params,
			)
		}
	}
}
