package jsonapi

import (
	"io/ioutil"
	"net/http"
)

// NewRequest ...
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

// Request ...
type Request struct {
	Method string
	URL    *URL
	Doc    *Document
	User   string
}
