package jsonapi

import "encoding/json"

type payloadSkeleton struct {
	Data     json.RawMessage        `json:"data"`
	Included []json.RawMessage      `json:"included"`
	Meta     map[string]interface{} `json:"meta"`
}
