package jsonapi

// Collection ...
type Collection interface {
	Type() string
	Len() int
	Elem(i int) Resource
	Add(r Resource)

	// JSON
	UnmarshalJSON(payload []byte) error
}
