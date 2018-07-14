package jsonapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error represents an error object from the JSON API specification.
type Error struct {
	ID     string                 `json:"id"`
	Code   string                 `json:"code"`
	Status int                    `json:"status"`
	Title  string                 `json:"title"`
	Detail string                 `json:"detail"`
	Links  map[string]string      `json:"links"`
	Source map[string]interface{} `json:"source"`
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

// NewErrBadRequest (400) ...
func NewErrBadRequest(title, detail string) Error {
	return Error{
		Status: http.StatusBadRequest,
		Title:  title,
		Detail: detail,
	}
}

// NewErrUnauthorized (401) ...
func NewErrUnauthorized() Error {
	return Error{
		Status: http.StatusUnauthorized,
		Title:  "Unauthorized",
		Detail: "Authentification is required to perform this request.",
	}
}

// NewErrForbidden (403) ...
func NewErrForbidden() Error {
	return Error{
		Status: http.StatusForbidden,
		Title:  "Forbidden",
		Detail: "Permission is required to perform this request.",
	}
}

// NewErrNotFound (404) ...
func NewErrNotFound() Error {
	return Error{
		Status: http.StatusNotFound,
		Title:  "Not found",
		Detail: "The URI does not exist.",
	}
}

// NewErrPayloadTooLarge (413) ...
func NewErrPayloadTooLarge() Error {
	return Error{
		Status: http.StatusRequestEntityTooLarge,
		Title:  "Payload too large",
		Detail: "That's what she said.",
	}
}

// NewErrRequestURITooLong (414) ...
func NewErrRequestURITooLong() Error {
	return Error{
		Status: http.StatusRequestURITooLong,
		Title:  "URI too long",
	}
}

// NewErrUnsupportedMediaType (415) ...
func NewErrUnsupportedMediaType() Error {
	return Error{
		Status: http.StatusUnsupportedMediaType,
		Title:  "Unsupported media type",
	}
}

// NewErrTooManyRequests (429) ...
func NewErrTooManyRequests() Error {
	return Error{
		Status: http.StatusTooManyRequests,
		Title:  "Too many requests",
	}
}

// NewErrRequestHeaderFieldsTooLarge (431) ...
func NewErrRequestHeaderFieldsTooLarge() Error {
	return Error{
		Status: http.StatusRequestHeaderFieldsTooLarge,
		Title:  "Header fields too large",
	}
}

// NewErrInternalServerError (500) ...
func NewErrInternalServerError() Error {
	return Error{
		Status: http.StatusInternalServerError,
		Title:  "Internal server error",
	}
}

// NewErrServiceUnavailable (503) ...
func NewErrServiceUnavailable() Error {
	return Error{
		Status: http.StatusServiceUnavailable,
		Title:  "Service unavailable",
	}
}
