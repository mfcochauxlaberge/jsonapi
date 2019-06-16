package jsonapi_test

import (
	"net/http"
	"testing"

	. "github.com/mfcochauxlaberge/jsonapi"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name           string
		err            Error
		expectedString string
	}{
		{
			name: "empty",
			err: func() Error {
				e := NewError()
				return e
			}(),
			expectedString: "",
		}, {
			name: "title",
			err: func() Error {
				e := NewError()
				e.Title = "An error"
				return e
			}(),
			expectedString: "An error",
		}, {
			name: "detail",
			err: func() Error {
				e := NewError()
				e.Detail = "An error occurred."
				return e
			}(),
			expectedString: "An error occurred.",
		}, {
			name: "http status code",
			err: func() Error {
				e := NewError()
				e.Status = http.StatusInternalServerError
				return e
			}(),
			expectedString: "500 Internal Server Error",
		}, {
			name: "http status code and title",
			err: func() Error {
				e := NewError()
				e.Status = http.StatusInternalServerError
				e.Title = "Internal server error"
				return e
			}(),
			expectedString: "500 Internal Server Error: Internal server error",
		}, {
			name: "http status code and detail",
			err: func() Error {
				e := NewError()
				e.Status = http.StatusInternalServerError
				e.Detail = "An internal server error occurred."
				return e
			}(),
			expectedString: "500 Internal Server Error: An internal server error occurred.",
		},
	}

	for _, test := range tests {
		assert.Equal(test.err.Error(), test.expectedString, test.name)
	}
}
