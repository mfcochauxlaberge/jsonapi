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

// NewError returns an empty Error object.
func NewError() Error {
	err := Error{
		ID:     "",
		Code:   "",
		Status: 0,
		Title:  "",
		Detail: "",
		Links:  map[string]string{},
		Source: map[string]interface{}{},
		Meta:   map[string]interface{}{},
	}

	return err
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
	e := NewError()

	e.Status = http.StatusBadRequest
	e.Title = title
	e.Detail = detail

	return e
}

// NewErrUnauthorized (401) ...
func NewErrUnauthorized() Error {
	e := NewError()

	e.Status = http.StatusUnauthorized
	e.Title = "Unauthorized"
	e.Detail = "Authentification is required to perform this request."

	return e
}

// NewErrForbidden (403) ...
func NewErrForbidden() Error {
	e := NewError()

	e.Status = http.StatusForbidden
	e.Title = "Forbidden"
	e.Detail = "Permission is required to perform this request."

	return e
}

// NewErrNotFound (404) ...
func NewErrNotFound() Error {
	e := NewError()

	e.Status = http.StatusNotFound
	e.Title = "Not found"
	e.Detail = "The URI does not exist."

	return e
}

// NewErrPayloadTooLarge (413) ...
func NewErrPayloadTooLarge() Error {
	e := NewError()

	e.Status = http.StatusRequestEntityTooLarge
	e.Title = "Payload too large"
	e.Detail = "That's what she said."

	return e
}

// NewErrRequestURITooLong (414) ...
func NewErrRequestURITooLong() Error {
	e := NewError()

	e.Status = http.StatusRequestURITooLong
	e.Title = "URI too long"

	return e
}

// NewErrUnsupportedMediaType (415) ...
func NewErrUnsupportedMediaType() Error {
	e := NewError()

	e.Status = http.StatusUnsupportedMediaType
	e.Title = "Unsupported media type"

	return e
}

// NewErrTooManyRequests (429) ...
func NewErrTooManyRequests() Error {
	e := NewError()

	e.Status = http.StatusTooManyRequests
	e.Title = "Too many requests"

	return e
}

// NewErrRequestHeaderFieldsTooLarge (431) ...
func NewErrRequestHeaderFieldsTooLarge() Error {
	e := NewError()

	e.Status = http.StatusRequestHeaderFieldsTooLarge
	e.Title = "Header fields too large"

	return e
}

// NewErrInternalServerError (500) ...
func NewErrInternalServerError() Error {
	e := NewError()

	e.Status = http.StatusInternalServerError
	e.Title = "Internal server error"

	return e
}

// NewErrServiceUnavailable (503) ...
func NewErrServiceUnavailable() Error {
	e := NewError()

	e.Status = http.StatusServiceUnavailable
	e.Title = "Service unavailable"

	return e
}
