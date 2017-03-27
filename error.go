package jsonapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

// MarshalJSON ...
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"status": strconv.FormatInt(int64(e.Status), 10),
		"title":  e.Title,
		"detail": e.Detail,
	})
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
