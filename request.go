package jsonapi

import (
	"io/ioutil"
	"net/http"
)

// NewRequest builds and returns a *Request based on r and schema.
//
// schema can be nil, in which case no checks will be done to insure that the
// request respects a specific schema.
func NewRequest(r *http.Request, schema *Schema) (*Request, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	url, err := NewURLFromRaw(schema, r.URL.EscapedPath())
	if err != nil {
		return nil, err
	}

	var doc *Document

	if r.Method == http.MethodPost || r.Method == http.MethodPatch {
		doc, err = UnmarshalDocument(body, schema)
		if err != nil {
			return nil, err
		}
	}

	req := &Request{
		Method: r.Method,
		URL:    url,
		Doc:    doc,
	}

	return req, nil
}

// A Request represents a JSON:API request.
type Request struct {
	Method string
	URL    *URL
	Doc    *Document
}
