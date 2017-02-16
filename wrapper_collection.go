package jsonapi

import "encoding/json"

// WrapCollection ...
func WrapCollection(r Resource) *WrapperCollection {
	// if r2, ok := v.(Resource); ok {
	// 	r = r2
	// } else {
	// 	r := Wrap(v)
	// }

	return &WrapperCollection{
		col:    []*Wrapper{},
		sample: r,
	}
}

// WrapperCollection ...
type WrapperCollection struct {
	col    []*Wrapper
	sample Resource
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

// Sample ...
func (wc *WrapperCollection) Sample() Resource {
	return wc.sample.New()
}

// MarshalJSONParams ...
func (wc *WrapperCollection) MarshalJSONParams(params *Params) ([]byte, error) {
	var raws []*json.RawMessage

	for _, r := range wc.col {
		var raw json.RawMessage
		raw, err := r.MarshalJSONParams(params)
		if err != nil {
			return []byte{}, err
		}
		raws = append(raws, &raw)
	}

	return json.Marshal(raws)
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
