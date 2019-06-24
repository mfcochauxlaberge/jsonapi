package jsonapi

import "encoding/json"

var _ Collection = (*WrapperCollection)(nil)

// WrapCollection ...
func WrapCollection(r Resource) *WrapperCollection {
	// if r2, ok := v.(Resource); ok {
	// 	r = r2
	// } else {
	// 	r := Wrap(v)
	// }

	typ := r.GetType().Name

	return &WrapperCollection{
		typ:    typ,
		col:    []*Wrapper{},
		sample: r,
	}
}

// WrapperCollection ...
type WrapperCollection struct {
	typ    string
	col    []*Wrapper
	sample Resource
}

// Type ....
func (wc *WrapperCollection) Type() string {
	return wc.typ
}

// Len ...
func (wc *WrapperCollection) Len() int {
	return len(wc.col)
}

// Elem ...
func (wc *WrapperCollection) Elem(i int) Resource {
	if len(wc.col) > i {
		return wc.col[i]
	}

	return nil
}

// Add ...
func (wc *WrapperCollection) Add(r Resource) {
	if wr, ok := r.(*Wrapper); ok {
		wc.col = append(wc.col, wr)
	}
}

// UnmarshalJSON ...
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
