package jsonapi

import "encoding/json"

type documentSkeleton struct {
	Data     json.RawMessage            `json:"data"`
	Included []json.RawMessage          `json:"included"`
	Links    map[string]json.RawMessage `json:"links"`
	Meta     map[string]interface{}     `json:"meta"`
	JSONAPI  map[string]interface{}     `json:"jsonapi"`
	Errors   []json.RawMessage          `json:"errors"`
}

type resourceSkeleton struct {
	ID            string                          `json:"id"`
	Attributes    json.RawMessage                 `json:"attributes"`
	Relationships map[string]relationshipSkeleton `json:"relationships"`
}

type identifierSkeleton struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type relationshipSkeleton struct {
	Links linksSkeleton
	Data  json.RawMessage
}

type linksSkeleton struct {
	Self    json.RawMessage
	Related json.RawMessage
}
