package jsonapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSelfLink(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name           string
		res            Resource
		id             string
		expectedString string
	}{
		{
			name:           "simple resource url",
			id:             "1",
			expectedString: "http://example.com/type/1",
		}, {
			name:           "simple resource url with hyphen in id",
			id:             "abc-123",
			expectedString: "http://example.com/type/abc-123",
		}, {
			name:           "empty id",
			id:             "",
			expectedString: "",
		},
	}

	for _, test := range tests {
		res := &SoftResource{}
		res.SetType(&Type{Name: "type"})
		res.SetID(test.id)

		link := buildSelfLink(res, "http://example.com")
		assert.Equal(test.expectedString, link, test.name)
	}
}
