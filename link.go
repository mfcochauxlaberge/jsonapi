package jsonapi

import (
	"encoding/json"
)

// Link represents a JSON:API links object.
type Link struct {
	HRef string                 `json:"href"`
	Meta map[string]interface{} `json:"meta"`
}

// MarshalJSON builds the JSON representation of a Link object.
func (l Link) MarshalJSON() ([]byte, error) {
	if len(l.Meta) > 0 {
		return json.Marshal(l)
	}

	return json.Marshal(l.HRef)
}
