package jsonapi

import "encoding/json"

type payloadSkeleton struct {
	Data     json.RawMessage        `json:"data"`
	Included []json.RawMessage      `json:"included"`
	Meta     map[string]interface{} `json:"meta"`
}

type resourceSkeleton struct {
	ID            string                          `json:"id"`
	Type          string                          `json:"type"`
	Attributes    map[string]json.RawMessage      `json:"attributes"`
	Relationships map[string]relationshipSkeleton `json:"relationships"`
}

type relationshipSkeleton struct {
	Data  json.RawMessage
	Links map[string]json.RawMessage
	Meta  map[string]json.RawMessage
}
