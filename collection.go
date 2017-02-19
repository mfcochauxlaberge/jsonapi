package jsonapi

// Collection ...
type Collection interface {
	Elem(i int) Resource
	Add(r Resource)
	Sample() Resource

	// JSON
	Marshal(url *URL) ([]byte, error)
	UnmarshalJSON(payload []byte) error
}
