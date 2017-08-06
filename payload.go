package jsonapi

// Payload ...
type Payload struct {
	Data     interface{}
	Included map[string]Resource
	Meta     map[string]interface{}
	JSONAPI  map[string]interface{}
}

// NewPayload ...
func NewPayload(method string, url *URL, body []byte, r *Registry) (*Payload, error) {
	pl := &Payload{
		Included: map[string]Resource{},
		Meta:     map[string]interface{}{},
		JSONAPI:  map[string]interface{}{},
	}

	var err error
	if method == "POST" {
		if url.Type == "col" {
			// Create resource
			pl.Data = r.Resource(url.ResType)
			_, err = Unmarshal(body, pl.Data)
		} else if url.Type == "self" && url.IsCol {
			// Create relationships
			pl.Data = Identifiers{}
			_, err = Unmarshal(body, pl.Data)
		}
	} else if method == "PATCH" {
		if url.Type == "res" {
			// Update resource
			pl.Data = r.Resource(url.ResType)
			_, err = Unmarshal(body, pl.Data)
		} else if url.Type == "self" && url.IsCol {
			// Create relationships
			pl.Data = Identifiers{}
			_, err = Unmarshal(body, pl.Data)
		}
	}

	return pl, err
}
