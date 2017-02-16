package jsonapi

// An Extra struct contains all the side information found in a JSON API
// document. This means the jsonapi, meta, and links top-level members.
type Extra struct {
	Meta    map[string]interface{}
	JSONAPI map[string]interface{}
}
