package jsonapi

import (
	"testing"

	"kkaribu/tchek"
)

func TestBuildSelfLink(t *testing.T) {
	tests := []struct {
		res            Resource
		expectedString string
	}{
		{
			// 0
			res:            Wrap(&user{ID: "1"}),
			expectedString: "http://example.com/users/1",
		}, {
			// 1
			res:            Wrap(&book{ID: "abc-123"}),
			expectedString: "http://example.com/books/abc-123",
		}, {
			// 2
			res:            Wrap(&user{ID: ""}),
			expectedString: "",
		},
	}

	for i, test := range tests {
		link := buildSelfLink(test.res, "http://example.com/")
		tchek.AreEqual(t, i, test.expectedString, link)
	}
}
