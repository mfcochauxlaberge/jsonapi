package jsonapi_test

import (
	"net/url"
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestParseURL(t *testing.T) {
	assert := assert.New(t)

	// Schema
	schema := newMockSchema()

	tests := []struct {
		name          string
		url           string
		expectedURL   URL
		expectedError bool
	}{

		{
			name:          "empty",
			url:           ``,
			expectedURL:   URL{},
			expectedError: true,
		}, {
			name:          "empty path",
			url:           `https://example.com`,
			expectedURL:   URL{},
			expectedError: true,
		}, {
			name: "collection name only",
			url:  `mocktypes1`,
			expectedURL: URL{
				Fragments:       []string{"mocktypes1"},
				Route:           "/mocktypes1",
				BelongsToFilter: BelongsToFilter{},
				ResType:         "mocktypes1",
				ResID:           "",
				RelKind:         "",
				IsCol:           true,
			},
			expectedError: false,
		}, {
			name: "full url for collection",
			url:  `https://api.example.com/mocktypes1`,
			expectedURL: URL{
				Fragments:       []string{"mocktypes1"},
				Route:           "/mocktypes1",
				BelongsToFilter: BelongsToFilter{},
				ResType:         "mocktypes1",
				ResID:           "",
				RelKind:         "",
				IsCol:           true,
			},
			expectedError: false,
		}, {
			name: "full url for resource",
			url:  `https://example.com/mocktypes1/mc1-1`,
			expectedURL: URL{
				Fragments:       []string{"mocktypes1", "mc1-1"},
				Route:           "/mocktypes1/:id",
				BelongsToFilter: BelongsToFilter{},
				ResType:         "mocktypes1",
				ResID:           "mc1-1",
				RelKind:         "",
				IsCol:           false,
			},
			expectedError: false,
		}, {
			name: "full url for related relationship",
			url:  `https://example.com/mocktypes1/mc1-1/to-one`,
			expectedURL: URL{
				Fragments: []string{"mocktypes1", "mc1-1", "to-one"},
				Route:     "/mocktypes1/:id/to-one",
				BelongsToFilter: BelongsToFilter{
					Type: "mocktypes1",
					ID:   "mc1-1",
					Name: "to-one",
				},
				ResType: "mocktypes2",
				ResID:   "",
				RelKind: "related",
				Rel: Rel{
					Name:         "to-one",
					Type:         "mocktypes2",
					ToOne:        true,
					InverseName:  "",
					InverseType:  "mocktypes1",
					InverseToOne: false,
				},
				IsCol: false,
			},
			expectedError: false,
		}, {
			name: "111111",
			url:  `https://example.com/mocktypes1/mc1-1/relationships/to-many-from-one`,
			expectedURL: URL{
				Fragments: []string{"mocktypes1", "mc1-1", "relationships", "to-many-from-one"},
				Route:     "/mocktypes1/:id/relationships/to-many-from-one",
				BelongsToFilter: BelongsToFilter{
					Type:        "mocktypes1",
					ID:          "mc1-1",
					Name:        "to-many-from-one",
					InverseName: "to-one-from-many",
				},
				ResType: "mocktypes2",
				ResID:   "",
				RelKind: "self",
				Rel: Rel{
					Name:         "to-many-from-one",
					Type:         "mocktypes2",
					ToOne:        false,
					InverseName:  "to-one-from-many",
					InverseType:  "mocktypes1",
					InverseToOne: true,
				},
				IsCol: true,
			},
			expectedError: false,
		}, {
			name: "full url for self relationship",
			url:  `/mocktypes1/mc1-1/relationships/to-many-from-one`,
			expectedURL: URL{
				Fragments: []string{"mocktypes1", "mc1-1", "relationships", "to-many-from-one"},
				Route:     "/mocktypes1/:id/relationships/to-many-from-one",
				BelongsToFilter: BelongsToFilter{
					Type:        "mocktypes1",
					ID:          "mc1-1",
					Name:        "to-many-from-one",
					InverseName: "to-one-from-many",
				},
				ResType: "mocktypes2",
				ResID:   "",
				RelKind: "self",
				Rel: Rel{
					Name:         "to-many-from-one",
					Type:         "mocktypes2",
					ToOne:        false,
					InverseName:  "to-one-from-many",
					InverseType:  "mocktypes1",
					InverseToOne: true,
				},
				IsCol: true,
			},
			expectedError: false,
		}, {
			name: "full url for self relationship with params",
			url:  `/mocktypes1/mc1-1/relationships/to-many-from-one?fields[mocktypes2]=boolptr%2Cint8ptr`,
			expectedURL: URL{
				Fragments: []string{"mocktypes1", "mc1-1", "relationships", "to-many-from-one"},
				Route:     "/mocktypes1/:id/relationships/to-many-from-one",
				BelongsToFilter: BelongsToFilter{
					Type:        "mocktypes1",
					ID:          "mc1-1",
					Name:        "to-many-from-one",
					InverseName: "to-one-from-many",
				},
				ResType: "mocktypes2",
				ResID:   "",
				RelKind: "self",
				Rel: Rel{
					Name:         "to-many-from-one",
					Type:         "mocktypes2",
					ToOne:        false,
					InverseName:  "to-one-from-many",
					InverseType:  "mocktypes1",
					InverseToOne: true,
				},
				IsCol: true,
			},
			expectedError: false,
		},
	}

	for _, test := range tests {
		u, _ := url.Parse(makeOneLineNoSpaces(test.url))
		url, err := ParseRawURL(schema, u.String())
		assert.Equal(test.expectedError, err != nil, test.name)

		// test.expectedURL.Path = makeOneLineNoSpaces(test.expectedURL.Path)

		if !test.expectedError {
			url.Params = nil
			assert.Equal(test.expectedURL, *url, test.name)
		}
	}
}

func TestParseParams(t *testing.T) {
	assert := assert.New(t)

	// Schema
	schema := newMockSchema()
	mockTypes1, _ := schema.GetType("mocktypes1")
	mockTypes2, _ := schema.GetType("mocktypes2")

	tests := []struct {
		name           string
		url            string
		resType        string
		expectedParams Params
		expectedError  bool
	}{
		{
			name: "slash only",
			url:  `/`,
			expectedParams: Params{
				Fields:       map[string][]string{},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				Filter:       nil,
				SortingRules: []string{"id"},
				PageSize:     0,
				PageNumber:   0,
				Include:      [][]Rel{},
			},
			expectedError: false,
		}, {
			name: "question mark",
			url:  `?`,
			expectedParams: Params{
				Fields:       map[string][]string{},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				Filter:       nil,
				SortingRules: []string{"id"},
				PageSize:     0,
				PageNumber:   0,
				Include:      [][]Rel{},
			},
			expectedError: false,
		}, {
			name: "include, sort, pagination in multiple parts",
			url: `
				?include=
					to-many-from-one.to-one-from-many.to-one.to-many-from-many,
					to-one-from-one.to-many-from-many
				&sort=to-many,str,,-bool
				&page[number]=3
				&sort=uint8
				&include=
					to-many-from-one,
					to-many-from-many
				&page[size]=50
			`,
			resType: "mocktypes1",
			expectedParams: Params{
				Fields: map[string][]string{
					"mocktypes1": mockTypes1.Fields(),
					"mocktypes2": mockTypes2.Fields(),
				},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				Filter:       nil,
				SortingRules: []string{"str", "-bool", "uint8", "int", "int16", "int32", "int64", "int8", "time", "uint", "uint16", "uint32", "uint64", "id"},
				PageSize:     50,
				PageNumber:   3,
				Include: [][]Rel{
					{
						mockTypes1.Rels["to-many-from-many"],
					},
					{
						mockTypes1.Rels["to-many-from-one"],
						mockTypes2.Rels["to-one-from-many"],
						mockTypes1.Rels["to-one"],
						mockTypes2.Rels["to-many-from-many"],
					},
					{
						mockTypes1.Rels["to-one-from-one"],
						mockTypes2.Rels["to-many-from-many"],
					},
				},
			},
			expectedError: false,
		}, {
			name: "sort param with many escaped commas",
			url: `
				?include=
					to-many-from-one.to-one-from-many.to-one.to-many-from-many%2C
					to-one-from-one.to-many-from-many
				&sort=to-many%2Cstr,%2C%2C-bool
				&page[number]=3
				&sort=uint8
				&include=
					to-many-from-one,
					to-many-from-many
				&page[size]=50
			`,
			resType: "mocktypes1",
			expectedParams: Params{
				Fields: map[string][]string{
					"mocktypes1": mockTypes1.Fields(),
					"mocktypes2": mockTypes2.Fields(),
				},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				Filter:       nil,
				SortingRules: []string{"str", "-bool", "uint8", "int", "int16", "int32", "int64", "int8", "time", "uint", "uint16", "uint32", "uint64", "id"},
				PageSize:     50,
				PageNumber:   3,
				Include: [][]Rel{
					{
						mockTypes1.Rels["to-many-from-many"],
					},
					{
						mockTypes1.Rels["to-many-from-one"],
						mockTypes2.Rels["to-one-from-many"],
						mockTypes1.Rels["to-one"],
						mockTypes2.Rels["to-many-from-many"],
					},
					{
						mockTypes1.Rels["to-one-from-one"],
						mockTypes2.Rels["to-many-from-many"],
					},
				},
			},
			expectedError: false,
		}, {
			name: "sort param with many unescaped commas",
			url: `
				?include=
					to-many-from-one.to-one-from-many
				&sort=to-many,str,,,-bool
				&sort=uint8
				&include=
					to-many-from-many
					to-many-from-one,
				&page[number]=110
				&page[size]=90
			`,
			resType: "mocktypes1",
			expectedParams: Params{
				Fields: map[string][]string{
					"mocktypes1": mockTypes1.Fields(),
				},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				Filter:       nil,
				SortingRules: []string{"str", "-bool", "uint8", "int", "int16", "int32", "int64", "int8", "time", "uint", "uint16", "uint32", "uint64", "id"},
				PageSize:     90,
				PageNumber:   110,
				Include: [][]Rel{
					{
						mockTypes1.Rels["to-many-from-one"],
						mockTypes2.Rels["to-one-from-many"],
					},
				},
			},
			expectedError: false,
		}, {
			name: "filter label",
			url: `
				?filter=label
			`,
			resType: "mocktypes1",
			expectedParams: Params{
				Fields: map[string][]string{
					"mocktypes1": mockTypes1.Fields(),
				},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				FilterLabel:  "label",
				Filter:       nil,
				SortingRules: []string{"bool", "int", "int16", "int32", "int64", "int8", "str", "time", "uint", "uint16", "uint32", "uint64", "uint8", "id"},
				PageSize:     0,
				PageNumber:   0,
				Include:      [][]Rel{},
			},
			expectedError: false,
		},
	}

	for _, test := range tests {
		u, err := url.Parse(makeOneLineNoSpaces(test.url))
		assert.NoError(err, test.name)

		su, err := NewSimpleURL(u)
		assert.NoError(err, test.name)

		params, err := NewParams(schema, su, test.resType)
		assert.Equal(test.expectedError, err != nil, test.name)

		// Set Attrs and Rels
		for resType, fields := range test.expectedParams.Fields {
			for _, field := range fields {
				if res, ok := schema.GetType(resType); ok {
					if _, ok := res.Attrs[field]; ok {
						test.expectedParams.Attrs[resType] = append(test.expectedParams.Attrs[resType], res.Attrs[field])
					} else if typ, ok := schema.GetType(resType); ok {
						if _, ok := typ.Rels[field]; ok {
							test.expectedParams.Rels[resType] = append(test.expectedParams.Rels[resType], res.Rels[field])
						}
					}
				}
			}
		}

		if !test.expectedError {
			assert.Equal(test.expectedParams, *params, test.name)
		}
	}
}
