package jsonapi

import (
	"encoding/json"
	"strings"
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

// buildSelfLink builds a URL that points to the resource represented by the
// value v.
//
// prepath is prepended to the path and usually represents a scheme and a
// domain name.
func buildSelfLink(res Resource, prepath string) string {
	link := prepath

	if !strings.HasSuffix(prepath, "/") {
		link += "/"
	}

	id, _ := res.Get("id").(string)
	if id != "" && res.GetType().Name != "" {
		link += res.GetType().Name + "/" + id
	}

	return link
}

// buildRelationshipLinks builds a links object (according to the JSON:API
// specification) that include both the self and related members.
func buildRelationshipLinks(res Resource, prepath, rel string) map[string]string {
	return map[string]string{
		"self":    buildSelfLink(res, prepath) + "/relationships/" + rel,
		"related": buildSelfLink(res, prepath) + "/" + rel,
	}
}
