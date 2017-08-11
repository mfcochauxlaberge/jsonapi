package jsonapi

// Payload ...
type Payload struct {
	Data     interface{}
	Included map[string]Resource
	Meta     map[string]interface{}
	JSONAPI  map[string]interface{}
}
