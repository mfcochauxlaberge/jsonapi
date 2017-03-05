package jsonapi

import (
	"fmt"
	"net/http"
)

// Error represents an error object from the JSON API specification.
type Error struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s: %s", e.Status, e.Title, e.Detail)
}

// NewErrInternal ...
func NewErrInternal() Error {
	return Error{
		Status: http.StatusInternalServerError,
		Title:  "Internet Server Error",
		Detail: "Something went wrong.",
	}
}

// NewErrNotFound ...
func NewErrNotFound() Error {
	return Error{
		Status: http.StatusNotFound,
		Title:  "Not Found",
		Detail: "The URI does not exist.",
	}
}

// NewErrBadRequest ...
func NewErrBadRequest() Error {
	return Error{
		Status: http.StatusNotFound,
		Title:  "Bad Request",
		Detail: "The content of the request is invalid.",
	}
}

// NewErrUnauthorized ...
func NewErrUnauthorized() Error {
	return Error{
		Status: http.StatusNotFound,
		Title:  "Unauthorized",
		Detail: "Identification is required to perform this request.",
	}
}

// NewErrForbidden ...
func NewErrForbidden() Error {
	return Error{
		Status: http.StatusNotFound,
		Title:  "Forbidden",
		Detail: "Permission is required to perform this request.",
	}
}
