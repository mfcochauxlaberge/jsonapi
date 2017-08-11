package jsonapi

import "encoding/json"

type payloadSkeleton struct {
	Data     json.RawMessage        `json:"data"`
	Included []json.RawMessage      `json:"included"`
	Meta     map[string]interface{} `json:"meta"`
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
