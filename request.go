package jsonapi

import (
	"encoding/json"
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

	if r.Method == http.MethodPatch || r.Method == http.MethodPost {
		doc = &Document{
			Included:  []Resource{},
			Resources: map[string]map[string]struct{}{},
			Links:     map[string]Link{},
			RelData:   map[string][]string{},
			Meta:      map[string]interface{}{},
		}
		ske := &payloadSkeleton{}

		// Unmarshal
		err = json.Unmarshal(body, ske)
		if err != nil {
			return nil, err
		}

		// Data
		if len(ske.Data) > 0 {
			if ske.Data[0] == '{' {
				// Resource
				res, err := UnmarshalResource(ske.Data, schema)
				if err != nil {
					return nil, err
				}
				doc.Data = res
			} else if ske.Data[0] == '[' {
				col, err := UnmarshalCollection(ske.Data, schema)
				if err != nil {
					return nil, err
				}
				doc.Data = col
			} else if string(ske.Data) == "null" {
				doc.Data = nil
			} else {
				// TODO Not exactly the right error
				return nil, NewErrMissingDataMember()
			}
		} else if len(ske.Errors) > 0 {
			doc.Errors = ske.Errors
		} else {
			return nil, NewErrMissingDataMember()
		}

		// Included
		if len(ske.Included) > 0 {
			inc := Identifier{}
			incs := []Identifier{}
			for _, rawInc := range ske.Included {
				err = json.Unmarshal(rawInc, &inc)
				if err != nil {
					return nil, err
				}
				incs = append(incs, inc)
			}

			for i := range incs {
				res, err := UnmarshalResource(ske.Included[i], schema)
				if err != nil {
					return nil, err
				}
				doc.Included = append(doc.Included, res)
			}
		}

		// Meta
		doc.Meta = ske.Meta
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
