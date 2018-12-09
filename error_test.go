package jsonapi

import (
	"net/http"
	"testing"

	"github.com/kkaribu/tchek"
)

func TestError(t *testing.T) {
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
				e.Detail = "An error occured."
				return e
			}(),
			expectedString: "An error occured.",
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
				e.Detail = "An internal server error occured."
				return e
			}(),
			expectedString: "500 Internal Server Error: An internal server error occured.",
		},
	}

	for _, test := range tests {
		tchek.AreEqual(t, test.name, test.err.Error(), test.expectedString)
	}
}
