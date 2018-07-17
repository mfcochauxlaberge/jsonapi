package jsonapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/twinj/uuid"
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
		ID:     uuid.NewV4().String(),
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
		if e.Detail != "" {
			return fmt.Sprintf("%d %s: %s", e.Status, fullName, e.Detail)
		} else if e.Title != "" {
			return fmt.Sprintf("%d %s: %s", e.Status, fullName, e.Title)
		} else {
			return fmt.Sprintf("%d %s", e.Status, fullName)
		}
	}

	if e.Detail != "" {
		return e.Detail
	} else {
		return e.Title
	}
}

// MarshalJSON ...
func (e Error) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}

	if e.ID != "" {
		m["id"] = e.ID
	}

	if e.Code != "" {
		m["code"] = e.Code
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

	if len(e.Links) > 0 {
		m["links"] = e.Links
	}

	if len(e.Source) > 0 {
		m["source"] = e.Source
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

// NewErrMalformedFilterParameter (400) ...
func NewErrMalformedFilterParameter(badFitler string) Error {
	e := NewError()

	e.Status = http.StatusBadRequest
	e.Title = "Malformed filter parameter"
	e.Detail = "The filter parameter is not a valid JSON object."
	e.Source["parameter"] = "filter"
	e.Meta["bad-filter"] = badFitler

	return e
}

// NewErrInvalidPageNumberParameter (400) ...
func NewErrInvalidPageNumberParameter(badPageNumber string) Error {
	e := NewError()

	e.Status = http.StatusBadRequest
	e.Title = "Invalid page number parameter"
	e.Detail = "The page number parameter is not positive integer (including 0)."
	e.Source["parameter"] = "page[number]"
	e.Meta["bad-page-number"] = badPageNumber

	return e
}

// NewErrInvalidPageSizeParameter (400) ...
func NewErrInvalidPageSizeParameter(badPageSize string) Error {
	e := NewError()

	e.Status = http.StatusBadRequest
	e.Title = "Invalid page size parameter"
	e.Detail = "The page size parameter is not positive integer (including 0)."
	e.Source["parameter"] = "page[size]"
	e.Meta["bad-page-size"] = badPageSize

	return e
}

// NewErrUnknownParameter (400) ...
func NewErrUnknownParameter(param string) Error {
	e := NewError()

	e.Status = http.StatusBadRequest
	e.Title = "Unknown parameter"
	e.Detail = fmt.Sprintf("%s is not a known parameter.", param)
	e.Source["parameter"] = param
	e.Meta["unknown-parameter"] = param

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
