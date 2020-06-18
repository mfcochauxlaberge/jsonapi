package jsonapi

// Copier ...
type Copier interface {
	New() Resource
	Copy() Resource
}
