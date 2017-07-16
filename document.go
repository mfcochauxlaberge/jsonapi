package jsonapi

import (
	"encoding/json"
)

// Document ...
type Document struct {
	// Data
	Resource    Resource
	Collection  Collection
	Identifier  Identifier
	Identifiers Identifiers

	// Included
	Included map[string]Resource

	// References
	Resources map[string]map[string]struct{}
	Links     map[string]Link

	// Relationships where data has to be included in payload
	RelData map[string][]string

	// Top-level members
	Meta    map[string]interface{}
	JSONAPI map[string]interface{}

	// Errors
	Errors []Error

	// URL
	URL *URL
}

// NewDocument ...
func NewDocument() *Document {
	return &Document{
		Included:  map[string]Resource{},
		Resources: map[string]map[string]struct{}{},
		Links:     map[string]Link{},
		RelData:   map[string][]string{},
		URL:       NewURL(),
	}
}

// Include ...
func (d *Document) Include(res Resource) {
	id, typ := res.IDAndType()
	key := typ + " " + id

	if len(d.Included) == 0 {
		d.Included = map[string]Resource{}
	}

	// Check resource
	if d.Resource != nil {
		rid, rtype := d.Resource.IDAndType()
		rkey := rid + " " + rtype

		if rkey == key {
			return
		}
	}

	// Check Collection
	if d.Collection != nil {
		_, ctyp := d.Collection.Sample().IDAndType()
		if ctyp == typ {
			for i := 0; i < d.Collection.Len(); i++ {
				rid, rtype := d.Collection.Elem(i).IDAndType()
				rkey := rid + " " + rtype

				if rkey == key {
					return
				}
			}
		}
	}

	// Check already included resources
	if _, ok := d.Included[key]; ok {
		return
	}

	d.Included[key] = res
}

// MarshalJSON ...
func (d *Document) MarshalJSON() ([]byte, error) {
	// Data
	var data json.RawMessage
	var errors json.RawMessage
	var err error
	if d.Errors != nil {
		errors, err = json.Marshal(d.Errors)
	} else if d.Resource != nil {
		_, typ := d.Resource.IDAndType()
		data, err = marshalResource(d.Resource, d.URL.Host, d.URL.Params.Fields[typ], d.RelData)
	} else if d.Collection != nil {
		data, err = marshalCollection(d.Collection, d.URL.Host, d.URL.Params.Fields[d.Collection.Type()], d.RelData)
	} else if (d.Identifier != Identifier{}) {
		data, err = json.Marshal(d.Identifier)
	} else if d.Identifiers != nil {
		data, err = json.Marshal(d.Identifiers)
	} else {
		data = []byte("null")
	}

	if err != nil {
		return []byte{}, err
	}

	// Included
	inclusions := []*json.RawMessage{}
	if len(data) > 0 {
		for key := range d.Included {
			_, typ := d.Included[key].IDAndType()
			raw, err := marshalResource(d.Included[key], d.URL.Host, d.URL.Params.Fields[typ], d.RelData)
			if err != nil {
				return []byte{}, err
			}
			rawm := json.RawMessage(raw)
			inclusions = append(inclusions, &rawm)
		}
	}

	// Marshaling
	plMap := map[string]interface{}{}

	if len(data) > 0 {
		plMap["data"] = data
	}

	if len(d.Links) > 0 {
		plMap["links"] = d.Links
	}

	if len(errors) > 0 {
		plMap["errors"] = errors
	}

	if len(inclusions) > 0 {
		plMap["included"] = inclusions
	}

	if len(d.Meta) > 0 {
		plMap["meta"] = d.Meta
	}

	if len(d.JSONAPI) > 0 {
		plMap["jsonapi"] = d.JSONAPI
	}

	return json.Marshal(plMap)
}
