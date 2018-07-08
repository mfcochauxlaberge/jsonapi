package jsonapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error represents an error object from the JSON API specification.
type Error struct {
	ID     string                 `json:"id"`
	Status int                    `json:"status"`
	Title  string                 `json:"title"`
	Detail string                 `json:"detail"`
	Meta   map[string]interface{} `json:"meta"`
}

// Error returns the string representation of the error.
//
// If the error does note contain a valid error status code, it returns an
// empty string.
func (e Error) Error() string {
	fullName := http.StatusText(e.Status)

	if fullName != "" && e.Status >= 400 && e.Status <= 599 {
		return fmt.Sprintf("%d %s: %s", e.Status, fullName, e.Title)
	}

	return ""
}

// MarshalJSON ...
func (e Error) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}

	if e.ID != "" {
		m["id"] = e.ID
	}

	if e.Status >= 400 && e.Status <= 599 {
		m["status"] = fmt.Sprintf("%d", e.Status)
	}

	if e.Title != "" {
		m["title"] = e.Title
	}

	if e.Detail != "" {
		m["detail"] = e.Detail
	}

	if len(e.Meta) > 0 {
		m["meta"] = e.Meta
	}

	return json.Marshal(m)
}

// NewErrInternal ...
func NewErrInternal() Error {
	return Error{
		Status: http.StatusInternalServerError,
		Title:  "Internal server error",
	}
}

// NewErrServiceUnavailable ...
func NewErrServiceUnavailable() Error {
	return Error{
		Status: http.StatusServiceUnavailable,
		Title:  "Service unavailable",
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
		Status: http.StatusBadRequest,
		Title:  "Bad Request",
		Detail: "One or more attributes are invalid.",
	}
}

// NewErrInvalidField ...
func NewErrInvalidField(detail string) Error {
	return Error{
		Status: http.StatusBadRequest,
		Title:  "Invalid Attribute",
		Detail: detail,
	}
}

// NewErrUnauthorized ...
func NewErrUnauthorized() Error {
	return Error{
		Status: http.StatusUnauthorized,
		Title:  "Unauthorized",
		Detail: "Identification is required to perform this request.",
	}
}

// NewErrForbidden ...
func NewErrForbidden() Error {
	return Error{
		Status: http.StatusForbidden,
		Title:  "Forbidden",
		Detail: "Permission is required to perform this request.",
	}
}
