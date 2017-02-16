package jsonapi

import "encoding/json"

// Document ...
type Document struct {
	// Data
	Resource    Resource
	Collection  Collection
	Identifier  Identifier
	Identifiers Identifiers

	// Included
	Included map[string]map[string]Resource

	// References
	Resources map[string]map[string]struct{}
	Linkage   map[string]map[string]struct{}

	// Extra
	// Links   map[string]interface{}
	Meta    map[string]interface{}
	JSONAPI map[string]interface{}

	// Errors
	Errors []Error

	// Params
	Params *Params
}

// Include ...
func (d *Document) Include(v interface{}) error {
	return nil
}

// MarshalJSON ...
func (d *Document) MarshalJSON() ([]byte, error) {
	// Data
	var data json.RawMessage
	var err error
	if d.Resource != nil {
		data, err = d.Resource.MarshalJSONParams(d.Params)
	} else if d.Collection != nil {
		data, err = d.Collection.MarshalJSONParams(d.Params)
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
	inclusions := []string{}
	for t := range d.Included {
		for id := range d.Included[t] {
			raw, err := json.Marshal(d.Included[t][id])
			if err != nil {
				return []byte{}, err
			}
			inclusions = append(inclusions, string(raw))
		}
	}

	// Marshaling
	plMap := map[string]interface{}{
		"data": &data,
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
