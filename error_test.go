package jsonapi

import (
	"net/http"
	"testing"

	"github.com/kkaribu/tchek"
)

func TestError(t *testing.T) {
	tests := []struct {
		err            Error
		expectedString string
	}{
		{
			// 0
			err: func() Error {
				e := NewError()
				return e
			}(),
			expectedString: "",
		}, {
			// 1
			err: func() Error {
				e := NewError()
				e.Title = "An error"
				return e
			}(),
			expectedString: "An error",
		}, {
			// 2
			err: func() Error {
				e := NewError()
				e.Detail = "An error occured."
				return e
			}(),
			expectedString: "An error occured.",
		}, {
			// 3
			err: func() Error {
				e := NewError()
				e.Status = http.StatusInternalServerError
				return e
			}(),
			expectedString: "500 Internal Server Error",
		}, {
			// 4
			err: func() Error {
				e := NewError()
				e.Status = http.StatusInternalServerError
				e.Title = "Internal server error"
				return e
			}(),
			expectedString: "500 Internal Server Error: Internal server error",
		}, {
			// 5
			err: func() Error {
				e := NewError()
				e.Status = http.StatusInternalServerError
				e.Detail = "An internal server error occured."
				return e
			}(),
			expectedString: "500 Internal Server Error: An internal server error occured.",
		},
	}

	for n, test := range tests {
		tchek.AreEqual(t, n, test.err.Error(), test.expectedString)
	}
}
