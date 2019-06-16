package jsonapi_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	assert := assert.New(t)

	// Schema
	schema := newMockSchema()

	tests := []struct {
		name          string
		method        string
		url           string
		schema        *Schema
		expectedError error
	}{
		{
			name:          "get collection (mock schema)",
			method:        "GET",
			url:           "/mocktypes1",
			schema:        schema,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		body := bytes.NewBufferString("")
		req := httptest.NewRequest(test.method, test.url, body)

		doc, err := NewRequest(req, test.schema)
		assert.Equal(test.expectedError, err, test.name)
		assert.Equal(test.method, doc.Method, test.name)
	}
}
