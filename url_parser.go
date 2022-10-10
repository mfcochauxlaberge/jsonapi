package jsonapi

import (
	"net/url"
	"strings"
)

// URLParser stores parsing options.
type URLParser struct {
	// Path prefix to ignore.
	//
	// Example: If PathPrefix is set to '/api/v1', then
	//          '/api/v1/objects' will become '/objects'.
	PathPrefix string
	Schema     *Schema
}

// Parse takes a raw path and returns a new one ready to be consumed by the
// library according to the given options.
func (o *URLParser) Parse(raw string) (*URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	parsed.Path = o.Path(parsed.Path)

	su, err := NewSimpleURL(parsed)
	if err != nil {
		return nil, err
	}

	u, err := NewURL(o.Schema, su)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Path takes a raw path and returns a new one ready to be consumed by the
// library according to the given options.
func (o *URLParser) Path(p string) string {
	path := strings.Trim(p, "/")
	prefix := strings.Trim(o.PathPrefix, "/")

	path = strings.TrimLeft(path, prefix)

	if len(path) > 0 && path[0] == '/' {
		return path
	}

	return "/" + path
}
