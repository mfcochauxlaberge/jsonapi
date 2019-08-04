package jsonapi

import "encoding/json"

var _ Collection = (*WrapperCollection)(nil)

// WrapCollection returns a *WrapperCollection which implements the Collection
// interface and holds resources of the type defined in r.
func WrapCollection(r Resource) *WrapperCollection {
	// if r2, ok := v.(Resource); ok {
	// 	r = r2
	// } else {
	// 	r := Wrap(v)
	// }

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

// UnmarshalJSON populates the receiver with the resources represented in the
// payload.
func (wc *WrapperCollection) UnmarshalJSON(payload []byte) error {
	var raws []json.RawMessage

	err := json.Unmarshal(payload, &raws)
	if err != nil {
		return err
	}

	for _, raw := range raws {
		r := wc.sample.New()
		err = json.Unmarshal(raw, r)
		if err != nil {
			wc.col = nil
			return err
		}
		wc.Add(r)
	}

	return nil
}
