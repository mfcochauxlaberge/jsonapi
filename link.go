package jsonapi

import "encoding/json"

// Link ...
type Link struct {
	HRef string                 `json:"href"`
	Meta map[string]interface{} `json:"meta"`
}

// MarshalJSON ...
func (l Link) MarshalJSON() ([]byte, error) {
	if len(l.Meta) > 0 {
		return json.Marshal(l)
	}

	return []byte(l.HRef), nil
}
