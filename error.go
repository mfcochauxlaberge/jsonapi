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

// NewInternalError ...
func NewInternalError() Error {
	return Error{
		Status: http.StatusInternalServerError,
		Title:  "Internet Server Error",
		Detail: "Something went wrong.",
	}
}

// NewNotFoundError ...
func NewNotFoundError() Error {
	return Error{
		Status: http.StatusNotFound,
		Title:  "Not Found",
		Detail: "The URI does not exist.",
	}
}

// NewBadRequestError ...
func NewBadRequestError() Error {
	return Error{
		Status: http.StatusNotFound,
		Title:  "Bad Request",
		Detail: "The content of the request is invalid.",
	}
}

// NewUnauthorizedError ...
func NewUnauthorizedError() Error {
	return Error{
		Status: http.StatusNotFound,
		Title:  "Unauthorized",
		Detail: "Identification is required to perform this request.",
	}
}

// NewForbiddenError ...
func NewForbiddenError() Error {
	return Error{
		Status: http.StatusNotFound,
		Title:  "Forbidden",
		Detail: "Permission is required to perform this request.",
	}
}

// Errors ...
type Errors []Error

// MarshalJSONParams ...
func (e Errors) MarshalJSONParams(_ Params) ([]byte, error) {
	return []byte{}, nil
}
