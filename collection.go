package jsonapi

// // A Collection defines the interface of a structure that can manage a set of
// // ordered resources of the same type.
// type Collection interface {
// 	// Type returns the name of the resources' type.
// 	GetType() Type

// 	// Len returns the number of resources in the collection.
// 	Len() int

// 	// At returns the resource at index i.
// 	At(int) Resource

// 	// Add adds a resource in the collection.
// 	Add(Resource)
// }

// Collection is a slice of objects that implement the Resource interface. They
// do not necessarily have the same type.
type Collection []Resource

// GetType returns a zero Type object because the collection does not represent
// a particular type.
func (r Collection) GetType() Type {
	return Type{}
}

// Len returns the number of elements in r.
func (r Collection) Len() int {
	return len(r)
}

// At returns the number of elements in r.
func (r Collection) At(i int) Resource {
	if i >= 0 && i < r.Len() {
		return (r)[i]
	}
	return nil
}

// Add adds a Resource object to r.
func (r *Collection) Add(res Resource) {
	*r = append(*r, res)
}
