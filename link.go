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
		var err error

		m := map[string]json.RawMessage{}

		m["href"], _ = json.Marshal(l.HRef)

		m["meta"], err = json.Marshal(l.Meta)
		if err != nil {
			return []byte{}, err
		}

		return json.Marshal(m)
	}

	return json.Marshal(l.HRef)
}
