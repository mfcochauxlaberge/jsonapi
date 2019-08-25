package jsonapi_test

import (
	"net/url"
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestParseURLPOC(t *testing.T) {
	assert := assert.New(t)

	// path := "/"

	assert.Equal(true, true, "obviously")
}

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
				IsCol:           true,
			},
			expectedError: false,
		}, {
			name:          "type not found",
			url:           "/mocktypes99",
			expectedError: true,
		}, {
			name:          "relationship not found",
			url:           "/mocktypes1/abc/relnotfound",
			expectedError: true,
		}, {
			name: "bad params",
			url: `
				/mocktypes1
				?fields[invalid]=attr1,attr2
			`,
			expectedError: true,
		}, {
			name:          "invalid raw url",
			url:           "%z",
			expectedError: true,
		}, {
			name: "invalid simpleurl",
			url: `
				/mocktypes1/abc123
				?page[size]=invalid
			`,
			expectedError: true,
		}, {
			name: "full url for collection",
			url:  `https://api.example.com/mocktypes1`,
			expectedURL: URL{
				Fragments:       []string{"mocktypes1"},
				Route:           "/mocktypes1",
				BelongsToFilter: BelongsToFilter{},
				ResType:         "mocktypes1",
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
				RelKind: "related",
				Rel: Rel{
					FromName: "to-one",
					FromType: "mocktypes1",
					ToOne:    true,
					ToName:   "",
					ToType:   "mocktypes2",
					FromOne:  false,
				},
			},
			expectedError: false,
		}, {
			name: "full url for self relationships",
			url: `
				https://example.com/mocktypes1/mc1-1/relationships/to-many-from-one
			`,
			expectedURL: URL{
				Fragments: []string{
					"mocktypes1", "mc1-1", "relationships", "to-many-from-one",
				},
				Route: "/mocktypes1/:id/relationships/to-many-from-one",
				BelongsToFilter: BelongsToFilter{
					Type:   "mocktypes1",
					ID:     "mc1-1",
					Name:   "to-many-from-one",
					ToName: "to-one-from-many",
				},
				ResType: "mocktypes2",
				RelKind: "self",
				Rel: Rel{
					FromName: "to-many-from-one",
					FromType: "mocktypes1",
					ToOne:    false,
					ToName:   "to-one-from-many",
					ToType:   "mocktypes2",
					FromOne:  true,
				},
				IsCol: true,
			},
			expectedError: false,
		}, {
			name: "path for self relationship",
			url:  `/mocktypes1/mc1-1/relationships/to-many-from-one`,
			expectedURL: URL{
				Fragments: []string{
					"mocktypes1", "mc1-1", "relationships", "to-many-from-one",
				},
				Route: "/mocktypes1/:id/relationships/to-many-from-one",
				BelongsToFilter: BelongsToFilter{
					Type:   "mocktypes1",
					ID:     "mc1-1",
					Name:   "to-many-from-one",
					ToName: "to-one-from-many",
				},
				ResType: "mocktypes2",
				RelKind: "self",
				Rel: Rel{
					FromName: "to-many-from-one",
					FromType: "mocktypes1",
					ToOne:    false,
					ToName:   "to-one-from-many",
					ToType:   "mocktypes2",
					FromOne:  true,
				},
				IsCol: true,
			},
			expectedError: false,
		}, {
			name: "path for self relationship with params",
			url: `
				/mocktypes1/mc1-1/relationships/to-many-from-one
				?fields[mocktypes2]=boolptr%2Cint8ptr
			`,
			expectedURL: URL{
				Fragments: []string{
					"mocktypes1", "mc1-1", "relationships", "to-many-from-one",
				},
				Route: "/mocktypes1/:id/relationships/to-many-from-one",
				BelongsToFilter: BelongsToFilter{
					Type:   "mocktypes1",
					ID:     "mc1-1",
					Name:   "to-many-from-one",
					ToName: "to-one-from-many",
				},
				ResType: "mocktypes2",
				RelKind: "self",
				Rel: Rel{
					FromName: "to-many-from-one",
					FromType: "mocktypes1",
					ToOne:    false,
					ToName:   "to-one-from-many",
					ToType:   "mocktypes2",
					FromOne:  true,
				},
				IsCol: true,
			},
			expectedError: false,
		},
	}

	for _, test := range tests {
		url, err := NewURLFromRaw(schema, makeOneLineNoSpaces(test.url))
		if test.expectedError {
			assert.Error(err)
		} else {
			assert.NoError(err)
			url.Params = nil
			assert.Equal(test.expectedURL, *url, test.name)
		}
	}
}

func TestParseParams(t *testing.T) {
	assert := assert.New(t)

	// Schema
	schema := newMockSchema()
	mockTypes1 := schema.GetType("mocktypes1")
	mockTypes2 := schema.GetType("mocktypes2")

	tests := []struct {
		name           string
		url            string
		colType        string
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
				SortingRules: []string{},
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
				SortingRules: []string{},
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
			colType: "mocktypes1",
			expectedParams: Params{
				Fields: map[string][]string{
					"mocktypes1": mockTypes1.Fields(),
					"mocktypes2": mockTypes2.Fields(),
				},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				SortingRules: []string{},
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
			colType: "mocktypes1",
			expectedParams: Params{
				Fields: map[string][]string{
					"mocktypes1": mockTypes1.Fields(),
					"mocktypes2": mockTypes2.Fields(),
				},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				SortingRules: []string{},
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
			colType: "mocktypes1",
			expectedParams: Params{
				Fields: map[string][]string{
					"mocktypes1": mockTypes1.Fields(),
				},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				SortingRules: []string{},
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
			colType: "mocktypes1",
			expectedParams: Params{
				Fields: map[string][]string{
					"mocktypes1": mockTypes1.Fields(),
				},
				Attrs:        map[string][]Attr{},
				Rels:         map[string][]Rel{},
				RelData:      map[string][]string{},
				FilterLabel:  "label",
				SortingRules: []string{},
				Include:      [][]Rel{},
			},
			expectedError: false,
		}, {
			name: "sorting rules without id",
			url: `
				/mocktypes1
				?sort=str,-int
			`,
			colType: "mocktypes1",
			expectedParams: Params{
				Fields: map[string][]string{
					"mocktypes1": mockTypes1.Fields(),
				},
				Attrs:   map[string][]Attr{},
				Rels:    map[string][]Rel{},
				RelData: map[string][]string{},
				SortingRules: []string{
					"str", "-int", "bool", "int16", "int32", "int64", "int8",
					"time", "uint", "uint16", "uint32", "uint64", "uint8", "id"},
				Include: [][]Rel{},
			},
			expectedError: false,
		}, {
			name: "sorting rules with id",
			url: `
				/mocktypes1
				?sort=str,-int,id
			`,
			colType: "mocktypes1",
			expectedParams: Params{
				Fields: map[string][]string{
					"mocktypes1": mockTypes1.Fields(),
				},
				Attrs:   map[string][]Attr{},
				Rels:    map[string][]Rel{},
				RelData: map[string][]string{},
				SortingRules: []string{
					"str", "-int", "id", "bool", "int16", "int32", "int64", "int8",
					"time", "uint", "uint16", "uint32", "uint64", "uint8"},
				Include: [][]Rel{},
			},
			expectedError: false,
		},
	}

	for _, test := range tests {
		u, err := url.Parse(makeOneLineNoSpaces(test.url))
		assert.NoError(err, test.name)

		su, err := NewSimpleURL(u)
		assert.NoError(err, test.name)

		params, err := NewParams(schema, su, test.colType)
		if test.expectedError {
			assert.Error(err, test.name)
		} else {
			assert.NoError(err, test.name)
		}

		// Set Attrs and Rels
		for colType, fields := range test.expectedParams.Fields {
			for _, field := range fields {
				if typ := schema.GetType(colType); typ.Name != "" {
					if _, ok := typ.Attrs[field]; ok {
						test.expectedParams.Attrs[colType] = append(
							test.expectedParams.Attrs[colType],
							typ.Attrs[field],
						)
					} else if typ := schema.GetType(colType); typ.Name != "" {
						if _, ok := typ.Rels[field]; ok {
							test.expectedParams.Rels[colType] = append(
								test.expectedParams.Rels[colType],
								typ.Rels[field],
							)
						}
					}
				}
			}
		}

		if test.expectedError {
			assert.Error(err, test.name)
		} else {
			assert.NoError(err, test.name)
			assert.Equal(test.expectedParams, *params, test.name)
		}
	}
}

func TestURLEscaping(t *testing.T) {
	assert := assert.New(t)

	schema := newMockSchema()

	tests := []struct {
		url       string
		escaped   string
		unescaped string
	}{
		{
			url: `
				/mocktypes1
				?fields[mocktypes1]=bool%2Cint8
				&page[number]=2
				&page[size]=10
				&filter=a_label
			`,
			escaped: `
				/mocktypes1
				?fields%5Bmocktypes1%5D=bool%2Cint8
				&filter=a_label
				&page%5Bnumber%5D=2
				&page%5Bsize%5D=10
				&sort=bool%2Cint%2Cint16%2Cint32%2Cint64%2Cint8%2Cstr%2Ctime%2C
				uint%2Cuint16%2Cuint32%2Cuint64%2Cuint8%2Cid
				`,
			unescaped: `
				/mocktypes1
				?fields[mocktypes1]=bool,int8
				&filter=a_label
				&page[number]=2
				&page[size]=10
				&sort=bool,int,int16,int32,int64,int8,str,time,uint,uint16,
					uint32,uint64,uint8,id
			`,
		},
	}

	for _, test := range tests {
		url, err := NewURLFromRaw(schema, makeOneLineNoSpaces(test.url))
		assert.NoError(err)
		assert.Equal(
			makeOneLineNoSpaces(test.escaped),
			url.String(),
		)
		assert.Equal(
			makeOneLineNoSpaces(test.unescaped),
			url.UnescapedString(),
		)
	}
}
