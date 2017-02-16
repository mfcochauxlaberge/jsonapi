package jsonapi

// Collection ...
type Collection interface {
	Elem(i int) Resource
	Add(r Resource)
	Sample() Resource

	// JSON
	MarshalJSONParams(params *Params) ([]byte, error)
	UnmarshalJSON(payload []byte) error
}
