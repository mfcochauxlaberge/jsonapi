package jsonapi

// Options ...
type Options struct {
	// Host for building links
	Host string

	// Fields to include in payload
	Fields map[string][]string

	// Relationships where data has to be included in payload
	RelData map[string][]string

	// Top-level members
	Meta    map[string]interface{}
	JSONAPI map[string]interface{}
}

// NewOptions ...
func NewOptions(host string, params *Params) *Options {
	if params == nil {
		params = &Params{
			Fields:  map[string][]string{},
			RelData: map[string][]string{},
		}
	}

	if params.Fields == nil {
		params.Fields = map[string][]string{}
	}

	if params.RelData == nil {
		params.RelData = map[string][]string{}
	}

	return &Options{
		Host:    host,
		Fields:  params.Fields,
		RelData: params.RelData,
		Meta:    map[string]interface{}{},
		JSONAPI: map[string]interface{}{},
	}
}
