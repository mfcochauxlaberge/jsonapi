package jsonapi

// Copier is a interface for objects that can return a new and empty instance or
// a deep copy of themselves.
type Copier interface {
	New() Resource
	Copy() Resource
}
