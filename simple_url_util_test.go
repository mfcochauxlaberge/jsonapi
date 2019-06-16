package jsonapi

import (
	"testing"

	"github.com/mfcochauxlaberge/tchek"
)

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
