package jsonapi

// A Collection defines the interface of a structure that can manage a set of
// ordered resources of the same type.
type Collection interface {
	// Type returns the name of the resources' type.
	GetType() Type

	// Len returns the number of resources in the collection.
	Len() int

	// At returns the resource at index i.
	At(int) Resource

	// Add adds a resource in the collection.
	Add(Resource)

	// UnmarshalJSON unmarshals the bytes that represent a collection
	// of resources into the struct that implements the interface.
	UnmarshalJSON([]byte) error
}
