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

func TestErrorConstructors(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		name     string
		err      Error
		expected string
	}{
		{
			name: "NewError",
			err: func() Error {
				e := NewError()
				return e
			}(),
			expected: "",
		}, {
			name: "NewErrBadRequest",
			err: func() Error {
				e := NewErrBadRequest("bad request", "error detail")
				return e
			}(),
			expected: "400 Bad Request: error detail",
		}, {
			name: "NewErrMalformedFilterParameter",
			err: func() Error {
				e := NewErrMalformedFilterParameter("filter")
				return e
			}(),
			expected: "400 Bad Request: " +
				"The filter parameter is not a string or a valid JSON object.",
		}, {
			name: "NewErrInvalidPageNumberParameter",
			err: func() Error {
				e := NewErrInvalidPageNumberParameter("9")
				return e
			}(),
			expected: "400 Bad Request: " +
				"The page number parameter is not positive integer (including 0).",
		}, {
			name: "NewErrInvalidPageSizeParameter",
			err: func() Error {
				e := NewErrInvalidPageSizeParameter("9")
				return e
			}(),
			expected: "400 Bad Request: " +
				"The page size parameter is not positive integer (including 0).",
		}, {
			name: "NewErrInvalidFieldValueInBody",
			err: func() Error {
				e := NewErrInvalidFieldValueInBody("field", "bad", "int")
				return e
			}(),
			expected: "400 Bad Request: " +
				"The field value is invalid for the expected type.",
		}, {
			name: "NewErrDuplicateFieldInFieldsParameter",
			err: func() Error {
				e := NewErrDuplicateFieldInFieldsParameter("type", "field")
				return e
			}(),
			expected: "400 Bad Request: " +
				"The fields parameter contains the same field more than once.",
		}, {
			name: "NewErrMissingDataMember",
			err: func() Error {
				e := NewErrMissingDataMember()
				return e
			}(),
			expected: "400 Bad Request: Missing data top-level member in payload.",
		}, {
			name: "NewErrUnknownFieldInBody",
			err: func() Error {
				e := NewErrUnknownFieldInBody("type", "field")
				return e
			}(),
			expected: "400 Bad Request: field is not a known field.",
		}, {
			name: "NewErrUnknownFieldInURL",
			err: func() Error {
				e := NewErrUnknownFieldInURL("field")
				return e
			}(),
			expected: "400 Bad Request: field is not a known field.",
		}, {
			name: "NewErrUnknownParameter",
			err: func() Error {
				e := NewErrUnknownParameter("param")
				return e
			}(),
			expected: "400 Bad Request: param is not a known parameter.",
		}, {
			name: "NewErrUnknownRelationshipInPath",
			err: func() Error {
				e := NewErrUnknownRelationshipInPath("type", "rel", "path")
				return e
			}(),
			expected: "400 Bad Request: rel is not a relationship of type.",
		}, {
			name: "NewErrUnknownTypeInURL",
			err: func() Error {
				e := NewErrUnknownTypeInURL("type")
				return e
			}(),
			expected: "400 Bad Request: type is not a known type.",
		}, {
			name: "NewErrUnknownFieldInFilterParameter",
			err: func() Error {
				e := NewErrUnknownFieldInFilterParameter("field")
				return e
			}(),
			expected: "400 Bad Request: field is not a known field.",
		}, {
			name: "NewErrUnknownOperatorInFilterParameter",
			err: func() Error {
				e := NewErrUnknownOperatorInFilterParameter("=>")
				return e
			}(),
			expected: "400 Bad Request: => is not a known operator.",
		}, {
			name: "NewErrInvalidValueInFilterParameter",
			err: func() Error {
				e := NewErrInvalidValueInFilterParameter("value", "string")
				return e
			}(),
			expected: "400 Bad Request: value is not a known value.",
		}, {
			name: "NewErrUnknownCollationInFilterParameter",
			err: func() Error {
				e := NewErrUnknownCollationInFilterParameter("collation")
				return e
			}(),
			expected: "400 Bad Request: collation is not a known collation.",
		}, {
			name: "NewErrUnknownFilterParameterLabel",
			err: func() Error {
				e := NewErrUnknownFilterParameterLabel("label")
				return e
			}(),
			expected: "400 Bad Request: label is not a known filter query label.",
		}, {
			name: "NewErrUnauthorized",
			err: func() Error {
				e := NewErrUnauthorized()
				return e
			}(),
			expected: "401 Unauthorized: Authentification is required to perform this request.",
		}, {
			name: "NewErrForbidden",
			err: func() Error {
				e := NewErrForbidden()
				return e
			}(),
			expected: "403 Forbidden: Permission is required to perform this request.",
		}, {
			name: "NewErrNotFound",
			err: func() Error {
				e := NewErrNotFound()
				return e
			}(),
			expected: "404 Not Found: The URI does not exist.",
		}, {
			name: "NewErrPayloadTooLarge",
			err: func() Error {
				e := NewErrPayloadTooLarge()
				return e
			}(),
			expected: "413 Request Entity Too Large: That's what she said.",
		}, {
			name: "NewErrRequestURITooLong",
			err: func() Error {
				e := NewErrRequestURITooLong()
				return e
			}(),
			expected: "414 Request URI Too Long: URI too long",
		}, {
			name: "NewErrUnsupportedMediaType",
			err: func() Error {
				e := NewErrUnsupportedMediaType()
				return e
			}(),
			expected: "415 Unsupported Media Type: Unsupported media type",
		}, {
			name: "NewErrTooManyRequests",
			err: func() Error {
				e := NewErrTooManyRequests()
				return e
			}(),
			expected: "429 Too Many Requests: Too many requests",
		}, {
			name: "NewErrRequestHeaderFieldsTooLarge",
			err: func() Error {
				e := NewErrRequestHeaderFieldsTooLarge()
				return e
			}(),
			expected: "431 Request Header Fields Too Large: Header fields too large",
		}, {
			name: "NewErrInternalServerError",
			err: func() Error {
				e := NewErrInternalServerError()
				return e
			}(),
			expected: "500 Internal Server Error: Internal server error",
		}, {
			name: "NewErrServiceUnavailable",
			err: func() Error {
				e := NewErrServiceUnavailable()
				return e
			}(),
			expected: "503 Service Unavailable: Service unavailable",
		}, {
			name: "NewErrNotImplemented",
			err: func() Error {
				e := NewErrNotImplemented()
				return e
			}(),
			expected: "501 Not Implemented: Not Implemented",
		},
	}

	for _, test := range tests {
		assert.Equal(test.expected, test.err.Error(), test.name)
	}
}
