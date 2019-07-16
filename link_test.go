package jsonapi_test

import (
	"testing"

	"github.com/mfcochauxlaberge/jsonapi"
	"github.com/stretchr/testify/assert"
)

func TestMarshalLink(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		link            jsonapi.Link
		expectedPayload string
	}{
		{
			link:            jsonapi.Link{},
			expectedPayload: `""`,
		}, {
			link: jsonapi.Link{
				HRef: "example.org",
			},
			expectedPayload: `"example.org"`,
		}, {
			link: jsonapi.Link{
				HRef: "example.org",
				Meta: map[string]interface{}{
					"s": "abc",
					"n": 123,
				},
			},
			expectedPayload: `{"href":"example.org","meta":{"n":123,"s":"abc"}}`,
		},
	}

	for _, test := range tests {
		pl, err := test.link.MarshalJSON()
		assert.NoError(err)
		assert.Equal(test.expectedPayload, string(pl))
	}
}
