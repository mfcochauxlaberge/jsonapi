package jsonapi

import (
	"testing"

	"github.com/kkaribu/tchek"
)

func TestBuildSelfLink(t *testing.T) {
	tests := []struct {
		name           string
		res            Resource
		expectedString string
	}{
		{
			name:           "simple resource url",
			res:            Wrap(&MockType1{ID: "1"}),
			expectedString: "http://example.com/mocktypes1/1",
		}, {
			name:           "simple resource url with hyphen in id",
			res:            Wrap(&MockType1{ID: "abc-123"}),
			expectedString: "http://example.com/mocktypes1/abc-123",
		}, {
			name:           "empty id",
			res:            Wrap(&MockType1{ID: ""}),
			expectedString: "",
		},
	}

	for _, test := range tests {
		link := buildSelfLink(test.res, "http://example.com")
		tchek.AreEqual(t, test.name, test.expectedString, link)
	}
}
