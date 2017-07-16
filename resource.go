package jsonapi

// Resource ...
type Resource interface {
	// Structure
	IDAndType() (string, string)
	Attrs() []Attr
	Rels() []Rel
	Attr(key string) Attr
	Rel(key string) Rel
	New() Resource

	// Read
	Get(key string) interface{}

	// Update
	SetID(id string)
	Set(key string, val interface{})

	// Read relationship
	GetToOne(key string) string
	GetToMany(key string) []string

	// Update relationship
	SetToOne(key string, rel string)
	SetToMany(key string, rels []string)

	// Validate
	Validate() []error

	// Copy
	Copy() Resource

	// JSON
	UnmarshalJSON(payload []byte) error
}
