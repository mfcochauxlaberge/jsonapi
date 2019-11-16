package jsonapi

// WrapCollection returns a *WrapperCollection which implements the Collection
// interface and holds resources of the type defined in r.
func WrapCollection(r Resource) *WrapperCollection {
	return &WrapperCollection{
		typ:    r.GetType(),
		col:    []*Wrapper{},
		sample: r,
	}
}

// WrapperCollection is a Collection of resources of a certain type defined
// using the WrapCollection constructor.
//
// Only resources of that type can be added to the collection and the type may
// not be modified.
type WrapperCollection struct {
	typ    Type
	col    []*Wrapper
	sample Resource
}

// GetType returns the type of the resources in the collection.
func (wc *WrapperCollection) GetType() Type {
	return wc.typ
}

// Len returns the number of elements in the collection.
func (wc *WrapperCollection) Len() int {
	return len(wc.col)
}

// At returns the resource at the given index.
//
// It returns nil if the index is greater than the number of resources in the
// collection.
func (wc *WrapperCollection) At(i int) Resource {
	if len(wc.col) > i {
		return wc.col[i]
	}

	return nil
}

// Add appends the given resource at the end of the collection.
func (wc *WrapperCollection) Add(r Resource) {
	if wr, ok := r.(*Wrapper); ok {
		wc.col = append(wc.col, wr)
	}
}
