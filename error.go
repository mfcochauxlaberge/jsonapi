package jsonapi

import "fmt"

// InternalError is the error returned when marshaling went wrong. In that
// case, the response should be a 500 Internal Error.
type InternalError struct {
}

func (i InternalError) Error() string {
	return "jsonapi: internal error while marshaling"
}

// Error represents an error object from the JSON API specification.
type Error struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%d %s: %s", e.Status, e.Title, e.Detail)
}

// Errors ...
type Errors []Error

// MarshalJSONParams ...
func (e Errors) MarshalJSONParams(_ Params) ([]byte, error) {
	return []byte{}, nil
}
