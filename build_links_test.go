package jsonapi

import (
	"testing"

	"github.com/kkaribu/tchek"
)

func TestBuildSelfLink(t *testing.T) {
	tests := []struct {
		res            Resource
		expectedString string
	}{
		{
			// 0
			res:            Wrap(&MockType1{ID: "1"}),
			expectedString: "http://example.com/mocktypes1/1",
		}, {
			// 1
			res:            Wrap(&MockType1{ID: "abc-123"}),
			expectedString: "http://example.com/mocktypes1/abc-123",
		}, {
			// 2
			res:            Wrap(&MockType1{ID: ""}),
			expectedString: "",
		},
	}

	for i, test := range tests {
		link := buildSelfLink(test.res, "http://example.com/")
		tchek.AreEqual(t, i, test.expectedString, link)
	}
}
