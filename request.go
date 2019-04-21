package jsonapi

import (
	"io/ioutil"
	"net/http"
)

// NewRequest builds a return a *Request based on r and schema.
//
// schema can be nil, in which case no checks will be done to insure that
// the request respects a specific schema.
func NewRequest(r *http.Request, schema *Schema) (*Request, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	url, err := ParseRawURL(schema, r.URL.RawPath)
	if err != nil {
		return nil, err
	}

	doc, err := Unmarshal(body, url, schema)
	if err != nil {
		return nil, err
	}

	req := &Request{
		Method: r.Method,
		URL:    url,
		Doc:    doc,
		User:   "",
	}

	return req, nil
}

// Request represents a JSON:API request.
type Request struct {
	Method string
	URL    *URL
	Doc    *Document
	User   string
}
